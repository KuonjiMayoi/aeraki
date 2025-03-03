BASEDIR=$(dirname "$0")/..

SCRIPTS_DIR=$BASEDIR/test/e2e/scripts
COMMON_DIR=$BASEDIR/test/e2e/common
export ISTIO_VERSION=1.10.0
export BUILD_TAG=latest

bash ${SCRIPTS_DIR}/istio.sh -y -f ${COMMON_DIR}/istio-config.yaml
bash ${SCRIPTS_DIR}/aeraki.sh

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/prometheus.yaml -n istio-system
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.10/samples/addons/grafana.yaml -n istio-system

kubectl create namespace kiali-operator
helm install \
    --set cr.create=true \
    --set cr.namespace=istio-system \
    --namespace kiali-operator \
    --repo https://kiali.org/helm-charts \
    kiali-operator \
    kiali-operator

kubectl apply -f $BASEDIR/demo/gateway/demo-ingress.yaml -n istio-system
kubectl apply -f $BASEDIR/demo/gateway/istio-ingressgateway.yaml -n istio-system

if [ $# == 0 ]
then
    bash $BASEDIR/demo/metaprotocol-dubbo/install.sh
    bash $BASEDIR/demo/metaprotocol-thrift/install.sh
    bash ${BASEDIR}/demo/kafka/install.sh
fi
