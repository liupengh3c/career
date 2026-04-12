package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ProtoToHiveTypeMap 将 protobuf 类型映射到 Hive 类型
var ProtoToHiveTypeMap = map[descriptorpb.FieldDescriptorProto_Type]string{
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:   "DOUBLE",  // 1
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT:    "FLOAT",   // 2
	descriptorpb.FieldDescriptorProto_TYPE_INT64:    "BIGINT",  // 3
	descriptorpb.FieldDescriptorProto_TYPE_UINT64:   "BIGINT",  // 4
	descriptorpb.FieldDescriptorProto_TYPE_INT32:    "INT",     // 5
	descriptorpb.FieldDescriptorProto_TYPE_FIXED64:  "BIGINT",  // 6
	descriptorpb.FieldDescriptorProto_TYPE_FIXED32:  "INT",     // 7
	descriptorpb.FieldDescriptorProto_TYPE_BOOL:     "BOOLEAN", // 8
	descriptorpb.FieldDescriptorProto_TYPE_STRING:   "STRING",  // 9
	descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "BINARY",  // 12
	descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "INT",     // 13
	descriptorpb.FieldDescriptorProto_TYPE_ENUM:     "INT",     // 14
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "INT",     // 15
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "BIGINT",  // 16
	descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "INT",     // 17
	descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "BIGINT",  // 18
}

// DescriptorMap 存储消息描述符的映射
type DescriptorMap map[string]*descriptorpb.DescriptorProto

// ParseDescFile 解析 .desc 文件
func ParseDescFile(filePath string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fds); err != nil {
		return nil, fmt.Errorf("解析 FileDescriptorSet 失败: %v", err)
	}

	return fds, nil
}

// RegisterMessages 递归注册所有消息到 DescriptorMap
func RegisterMessages(messages []*descriptorpb.DescriptorProto, pkg, parentPath string, dm DescriptorMap, verbose bool) {
	for _, msg := range messages {
		var fullName string
		if parentPath != "" {
			fullName = parentPath + "." + msg.GetName()
		} else if pkg != "" {
			fullName = pkg + "." + msg.GetName()
		} else {
			fullName = msg.GetName()
		}

		// 多键索引
		dm[fullName] = msg
		dm[msg.GetName()] = msg
		if strings.HasPrefix(fullName, ".") {
			dm[fullName[1:]] = msg
		}

		if verbose {
			fmt.Printf("%-40s %s\n", msg.GetName(),
				map[bool]string{true: "嵌套字段", false: "非嵌套字段"}[len(msg.GetNestedType()) > 0])
		}

		// 递归处理嵌套消息
		if len(msg.GetNestedType()) > 0 {
			RegisterMessages(msg.GetNestedType(), pkg, fullName, dm, verbose)
		}
	}
}

// LookupMessage 查找消息定义
func LookupMessage(typeName string, dm DescriptorMap, parentPath string) *descriptorpb.DescriptorProto {
	if msg, ok := dm[typeName]; ok {
		return msg
	}
	if msg, ok := dm[strings.TrimPrefix(typeName, ".")]; ok {
		return msg
	}
	if parentPath != "" {
		fullPath := parentPath + "." + typeName
		if msg, ok := dm[fullPath]; ok {
			return msg
		}
	}
	return nil
}

// MapProtoToHive 将 proto 字段映射为 Hive 类型字符串（用于 ARRAY 内部）
func MapProtoToHive(field *descriptorpb.FieldDescriptorProto, dm DescriptorMap, parentPath string) string {
	fieldType := field.GetType()

	var hiveType string

	// TYPE_MESSAGE = 11
	if fieldType == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		msgDesc := LookupMessage(field.GetTypeName(), dm, parentPath)
		if msgDesc != nil {
			// 构建 STRUCT 类型，递归处理子字段
			var structFields []string
			for _, subField := range msgDesc.GetField() {
				subType := MapProtoToHive(subField, dm, field.GetTypeName())
				structFields = append(structFields, fmt.Sprintf("%s:%s", subField.GetName(), subType))
			}
			hiveType = fmt.Sprintf("STRUCT<%s>", strings.Join(structFields, ","))
		} else {
			hiveType = "STRING"
		}
	} else if t, ok := ProtoToHiveTypeMap[fieldType]; ok {
		hiveType = t
	} else {
		hiveType = "STRING"
	}

	// 处理 repeated 字段（label: 3 = LABEL_REPEATED）
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		hiveType = fmt.Sprintf("ARRAY<%s>", hiveType)
	}

	return hiveType
}

// CollectExpandedColumns 递归展开消息字段，生成 Hive 建表列定义列表
func CollectExpandedColumns(msg *descriptorpb.DescriptorProto, dm DescriptorMap, prefix string) []string {
	var columns []string

	for _, field := range msg.GetField() {
		colName := prefix + field.GetName()
		if prefix == "" {
			colName = field.GetName()
		}
		isRepeated := field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
		fieldType := field.GetType()

		if isRepeated {
			// repeated 字段：保留为 ARRAY，停止展开
			switch fieldType {
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				msgDesc := LookupMessage(field.GetTypeName(), dm, msg.GetName())
				if msgDesc != nil {
					var structFields []string
					for _, sf := range msgDesc.GetField() {
						sfType := MapProtoToHive(sf, dm, field.GetTypeName())
						structFields = append(structFields, fmt.Sprintf("%s:%s", sf.GetName(), sfType))
					}
					inner := fmt.Sprintf("STRUCT<%s>", strings.Join(structFields, ","))
					columns = append(columns, fmt.Sprintf("%s ARRAY<%s>", colName, inner))
				} else {
					columns = append(columns, fmt.Sprintf("%s ARRAY<STRING>", colName))
				}
			default:
				scalarType := ProtoToHiveTypeMap[fieldType]
				if scalarType == "" {
					scalarType = "STRING"
				}
				columns = append(columns, fmt.Sprintf("%s ARRAY<%s>", colName, scalarType))
			}
		} else if fieldType == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// 非 repeated 自定义消息类型：递归展开
			msgDesc := LookupMessage(field.GetTypeName(), dm, msg.GetName())
			if msgDesc != nil {
				subPrefix := colName + "__"
				columns = append(columns, CollectExpandedColumns(msgDesc, dm, subPrefix)...)
			} else {
				columns = append(columns, fmt.Sprintf("%s STRING", colName))
			}
		} else {
			// 标量类型
			hiveType := ProtoToHiveTypeMap[fieldType]
			if hiveType == "" {
				hiveType = "STRING"
			}
			columns = append(columns, fmt.Sprintf("%s %s", colName, hiveType))
		}
	}

	return columns
}

// GenerateHiveCreateTable 从 proto descriptor 文件生成 Hive 创建表语句
func GenerateHiveCreateTable(descFile, messageName string, verbose bool) (string, error) {
	// 构建描述符映射
	dm := make(DescriptorMap)

	if verbose {
		fmt.Println("\n=== 注册消息 ===")
	}
	fds, err := ParseDescFile(descFile)
	if err != nil {
		return "", err
	}
	for _, file := range fds.GetFile() {
		pkg := file.GetPackage()
		RegisterMessages(file.GetMessageType(), pkg, "", dm, verbose)
	}
	messages, _ := json.Marshal(dm)
	fmt.Println(string(messages))
	// 提取表名
	tableName := strings.TrimSuffix(filepath.Base(descFile), filepath.Ext(descFile))
	var columns []string

	// 查找并处理指定的消息
	for _, file := range fds.GetFile() {
		for _, msg := range file.GetMessageType() {
			if messageName != "" && msg.GetName() != messageName {
				continue
			}

			// 使用展开逻辑生成列
			columns = append(columns, CollectExpandedColumns(msg, dm, "")...)

		}
	}

	// 升序排列列
	sort.Strings(columns)

	if verbose {
		// 打印列定义
		fmt.Printf("\n%s\n", strings.Repeat("=", 60))
		fmt.Printf("生成的 Hive 列定义（共 %d 列）\n", len(columns))
		fmt.Printf("%s\n", strings.Repeat("=", 60))
		for i, col := range columns {
			fmt.Printf("%3d. %s\n", i+1, col)
		}
		fmt.Printf("%s\n\n", strings.Repeat("=", 60))
	}

	// 生成 SQL 语句
	columnsStr := strings.Join(columns, ",\n    ")
	sqlTemplate := fmt.Sprintf(`
CREATE TABLE %s (
    %s
)
STORED AS PARQUET;
    `, tableName, columnsStr)

	return sqlTemplate, nil
}

func main() {
	// 获取脚本所在目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取可执行文件路径失败: %v", err)
	}
	// 对于开发环境，使用相对路径
	// descFile := filepath.Join(currentDir, "fm_pnc", "fm_pnc_tag.desc")
	descFile := filepath.Join(currentDir, "team.desc")

	// 如果文件不存在，尝试当前工作目录
	if _, err := os.Stat(descFile); os.IsNotExist(err) {
		descFile = "fm_pnc_tag.desc"
	}

	fmt.Printf("解析 desc 文件: %s\n", descFile)

	// sql, err := GenerateHiveCreateTable(descFile, "Tags", true)
	sql, err := GenerateHiveCreateTable(descFile, "Team", true)
	if err != nil {
		log.Fatalf("生成 Hive 建表语句失败: %v", err)
	}

	fmt.Println(sql)
}
