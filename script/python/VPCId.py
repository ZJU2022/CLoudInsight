import re

data = '''
{
    "Action": "DescribeVPCResponse",
    "DataSet": [
        {
            "CreateTime": 1657597522,
            "IPv6Network": "",
            "Name": "wyxVPC",
            "Network": [
                "192.168.0.0/19"
            ],
            "NetworkInfo": [
                {
                    "Network": "192.168.0.0/19",
                    "SubnetCount": 1
                }
            ],
            "OperatorName": "",
            "Remark": "",
            "SubnetCount": 1,
            "Tag": "Default",
            "TunnelId": 21081509,
            "UpdateTime": 1657597522,
            "VPCId": "uvnet-bs30zxfk",
            "VPCType": "DefinedVPC"
        },
        {            "CreateTime": 1609751848,
            "IPv6Network": "",
            "Name": "VPC",
            "Network": [
                "10.18.0.0/16"
            ],
            "NetworkInfo": [
                {
                    "Network": "10.18.0.0/16",
                    "SubnetCount": 1
                }
            ],
            "OperatorName": "",
            "Remark": "",
            "SubnetCount": 1,
            "Tag": "Default",
            "TunnelId": 19990487,
            "UpdateTime": 1609751850,
            "VPCId": "uvnet-y41odxvd",
            "VPCType": "DefinedVPC"
        },
        {
            "CreateTime": 1559618544,
            "IPv6Network": "",
            "Name": "DefaultVPC",
            "Network": [
                "10.9.0.0/16",
                "10.10.0.0/16",
                "10.19.0.0/16",
                "10.42.0.0/16"
            ],
            "NetworkInfo": [
                {
                    "Network": "10.9.0.0/16",
                    "SubnetCount": 1
                },
                {
                    "Network": "10.10.0.0/16",
                    "SubnetCount": 1
                },
                {
                    "Network": "10.19.0.0/16",
                    "SubnetCount": 1
                },
                {
                    "Network": "10.42.0.0/16",
                    "SubnetCount": 1
                }
            ],
            "OperatorName": "",
            "Remark": "",
            "SubnetCount": 4,
            "Tag": "Default",
            "TunnelId": 18458149,
            "UpdateTime": 1559618546,
            "VPCId": "uvnet-eevc2cxx",
            "VPCType": "DefaultVPC"
        }
    ],
    "RetCode": 0,
    "TotalCount": 10
}
'''

# 更灵活的正则表达式，处理空白字符和换行符
pattern = r'"Name":\s*"DefaultVPC".*?"VPCId":\s*"([a-zA-Z0-9]{14})"'

matches = re.findall(pattern, data, re.DOTALL)  # 使用 re.DOTALL 让 '.' 匹配换行符

if matches:
    for vpc_id in matches:
        print("VPCId:", vpc_id)
else:
    print("未找到匹配的 VPCId")
