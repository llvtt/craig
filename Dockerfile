FROM public.ecr.aws/lambda/provided:al2
ADD main /
ENTRYPOINT [ "/main" ]
