import argparse

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    # 添加参数
    parser.add_argument('pos', help='位置参数')
    parser.add_argument('-v', '--verbose', help='enable verbose mode',action="store_true")
    parser.add_argument('-b', '--batch', type=int, default=100, help='set batch size')
    parser.add_argument('-m', '--milestones', nargs='+', required=True, default=15)
    # 解析参数
    args = parser.parse_args()
    # 使用参数
    if args.verbose:
        print("verbose mode is on",args.verbose)
    print("pos参数:", args.pos)
    print("batch size:", args.batch)
    print("milestones:", args.milestones)