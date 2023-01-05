package assets

import (
	"embed"

	rabbitmqv1 "github.com/rabbitmq/cluster-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	//go:embed manifests/*
	manifests  embed.FS
	appsScheme = runtime.NewScheme()
	coreScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
	coreCodecs = serializer.NewCodecFactory(coreScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
	if err := corev1.AddToScheme(coreScheme); err != nil {
		panic(err)
	}
}

func GetDeploymentFromFile(name string) *appsv1.Deployment {
	deploymentBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}
	deploymentObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion),
		deploymentBytes,
	)
	if err != nil {
		panic(err)
	}
	return deploymentObject.(*appsv1.Deployment)
}

func GetServiceFromFile(name string) *corev1.Service {
	serviceBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}
	serviceObject, err := runtime.Decode(
		coreCodecs.UniversalDecoder(corev1.SchemeGroupVersion),
		serviceBytes,
	)
	if err != nil {
		panic(err)
	}
	return serviceObject.(*corev1.Service)
}

func GetIngressFromFile(name string) *networkingv1.Ingress {
	ingressBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}
	ingressObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(networkingv1.SchemeGroupVersion),
		ingressBytes,
	)
	if err != nil {
		panic(err)
	}
	return ingressObject.(*networkingv1.Ingress)
}

func GetRabbitmqClusterFromFile(name string) *rabbitmqv1.RabbitmqCluster {
	rabbitmqClusterBytes, err := manifests.ReadFile(name)
	if err != nil {
		panic(err)
	}
	rabbitmqClusterObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(rabbitmqv1.GroupVersion),
		rabbitmqClusterBytes,
	)
	if err != nil {
		panic(err)
	}
	return rabbitmqClusterObject.(*rabbitmqv1.RabbitmqCluster)
}
