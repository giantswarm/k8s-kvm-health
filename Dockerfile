FROM alpine:3.14.3

ADD ./k8s-kvm-health /k8s-kvm-health

ENTRYPOINT ["/k8s-kvm-health"]
CMD ["daemon"]
