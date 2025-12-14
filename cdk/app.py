#!/usr/bin/env python3
"""
AWS CDK 애플리케이션 진입점
MCP 서버 인프라를 정의합니다
"""
import aws_cdk as cdk
from stacks.mcp_lambda_stack import MCPLambdaStack

app = cdk.App()

# MCP Lambda Stack 배포
MCPLambdaStack(
    app,
    "mcp-server-stack",
    env=cdk.Environment(
        account="YOUR_ACCOUNT_ID",  # AWS 계정 ID로 변경 필요
        region="ap-northeast-2",
    ),
)

app.synth()
