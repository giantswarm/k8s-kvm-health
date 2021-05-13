FROM alpine:3.13.5

ADD ./k8s-kvm-health /k8s-kvm-health

ENTRYPOINT ["/k8s-kvm-health"]
CMD ["daemon"]
