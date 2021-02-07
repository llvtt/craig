FROM public.ecr.aws/lambda/provided:al2
ADD cloudwatch-events /
ADD slack-events /
ENTRYPOINT [ "/cloudwatch-events" ]
ENTRYPOINT [ "/slack-events" ]
