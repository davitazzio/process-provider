FROM gcr.io/distroless/static@sha256:a01d47d4036cae5a67a9619e3d06fa14a6811a2247b4da72b4233ece4efebd57

ARG TARGETOS
ARG TARGETARCH

COPY --from=busybox /bin/sh /bin/sh
ADD bin/$TARGETOS\_$TARGETARCH/provider /usr/local/bin/crossplane-processprovider-provider


USER 65532
ENTRYPOINT ["crossplane-processprovider-provider"]
