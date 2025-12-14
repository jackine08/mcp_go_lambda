"""
MCP Lambda Stack
Lambda 함수, API Gateway, CloudWatch 리소스 정의
"""
from aws_cdk import (
    aws_lambda as lambda_,
    aws_apigateway as apigw,
    aws_logs as logs,
    aws_iam as iam,
    aws_cloudwatch as cloudwatch,
    Stack,
    Duration,
    RemovalPolicy,
    CfnOutput,
)
from constructs import Construct
import os


class MCPLambdaStack(Stack):
    """MCP 서버 Lambda Stack"""

    def __init__(self, scope: Construct, id: str, **kwargs):
        super().__init__(scope, id, **kwargs)

        # ==================
        # Lambda 함수 정의
        # ==================
        lambda_fn = lambda_.Function(
            self,
            "MCPServerFunction",
            runtime=lambda_.Runtime.PROVIDED_AL2023,
            handler="bootstrap",
            code=lambda_.Code.from_asset(
                os.path.join(os.path.dirname(__file__), "..", "..", "lambda_dist")
            ),
            timeout=Duration.seconds(30),
            memory_size=256,
            environment={
                "ENVIRONMENT": "dev",
            },
            description="MCP Server Lambda Function",
        )

        # ==================
        # API Gateway (REST API)
        # ==================
        api = apigw.RestApi(
            self,
            "MCPApi",
            rest_api_name="mcp-api",
            description="MCP Server REST API",
            deploy_options=apigw.StageOptions(
                stage_name="dev",
                logging_level=apigw.MethodLoggingLevel.INFO,
                data_trace_enabled=True,
            ),
        )

        # /mcp 리소스
        mcp_resource = api.root.add_resource("mcp")

        # POST /mcp 메서드
        mcp_resource.add_method(
            "POST",
            apigw.LambdaIntegration(lambda_fn),
        )

        # ==================
        # CloudWatch 대시보드
        # ==================
        dashboard = cloudwatch.Dashboard(
            self,
            "MCPServerDashboard",
            dashboard_name="mcp-server-dev",
        )

        # Lambda 메트릭
        dashboard.add_widgets(
            cloudwatch.GraphWidget(
                title="Lambda Invocations",
                left=[
                    lambda_fn.metric_invocations(
                        statistic="Sum",
                    ),
                ],
            ),
            cloudwatch.GraphWidget(
                title="Lambda Duration",
                left=[
                    lambda_fn.metric_duration(
                        statistic="Average",
                    ),
                ],
            ),
            cloudwatch.GraphWidget(
                title="Lambda Errors",
                left=[
                    lambda_fn.metric_errors(
                        statistic="Sum",
                    ),
                ],
            ),
        )

        # ==================
        # CloudWatch 알람
        # ==================
        error_alarm = cloudwatch.Alarm(
            self,
            "MCPServerErrorAlarm",
            metric=lambda_fn.metric_errors(statistic="Sum"),
            threshold=1,
            evaluation_periods=1,
            alarm_description="Alert when Lambda function has errors",
            alarm_name="mcp-server-dev-errors",
        )

        # ==================
        # Outputs
        # ==================
        CfnOutput(
            self,
            "APIEndpoint",
            value=api.url_for_path("/mcp"),
            description="MCP API Endpoint",
            export_name="MCPApiEndpoint",
        )

        CfnOutput(
            self,
            "LambdaFunctionArn",
            value=lambda_fn.function_arn,
            description="Lambda Function ARN",
            export_name="MCPLambdaFunctionArn",
        )

        CfnOutput(
            self,
            "LambdaFunctionName",
            value=lambda_fn.function_name,
            description="Lambda Function Name",
            export_name="MCPLambdaFunctionName",
        )
