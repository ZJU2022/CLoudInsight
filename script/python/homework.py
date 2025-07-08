import re
from typing import List
from datetime import datetime, date
from typing import Dict, Tuple
import pygame
import time
import random


def read_file(filepath: str) -> str:
    """从指定路径读取文件内容
    
    Args:
        filepath: 文件路径
        
    Returns:
        文件内容字符串
    """
    try:
        with open(filepath, 'r', encoding='utf-8') as file:
            return file.read()
    except IOError as e:
        print(f"无法读取文件 {filepath}: {e}")
        return ""


def extract_default_vpc_resources(content: str) -> List[str]:
    """提取内容中所属VPC为DefaultVPC的资源ID
    
    Args:
        content: 要分析的文本内容
        
    Returns:
        匹配到的资源ID列表
    """
    # CSV格式中寻找VPC字段为DefaultVPC的资源ID
    pattern = r'"[^"]*","[^"]*","[^"]*","([^"]*)","[^"]*","DefaultVPC"'
    return re.findall(pattern, content)


def get_labor_day_holiday_2025() -> Dict[str, object]:
    """获取2025年中国五一劳动节放假安排
    
    Returns:
        包含放假日期和详情的字典
    """
    today = datetime.now().date()
    
    holiday_info = {
        "name": "劳动节", 
        "english_name": "Labor Day",
        "start_date": date(2025, 5, 1),
        "end_date": date(2025, 5, 5),
        "duration": 5,
        "description": "2025年劳动节放假安排为5月1日至5月5日，共5天。",
        "source": "中国政府网 (english.www.gov.cn)"
    }
    
    # 检查当前日期是否在假期内
    if today >= holiday_info["start_date"] and today <= holiday_info["end_date"]:
        holiday_info["status"] = "正在假期中"
    elif today < holiday_info["start_date"]:
        days_remaining = (holiday_info["start_date"] - today).days
        holiday_info["status"] = f"距离假期还有{days_remaining}天"
    else:
        holiday_info["status"] = "假期已结束"
    
    return holiday_info


def main() -> None:
    """主函数：读取文件，提取资源ID并打印结果"""
    input_file = 'mysql_list_20241119170205.txt'
    
    file_content = read_file(input_file)
    if not file_content:
        return
    
    resource_ids = extract_default_vpc_resources(file_content)
    
    if resource_ids:
        print("匹配到的资源ID:")
        for resource_id in resource_ids:
            print(resource_id)
    else:
        print("没有找到匹配的资源ID。")


if __name__ == "__main__":
    main()

