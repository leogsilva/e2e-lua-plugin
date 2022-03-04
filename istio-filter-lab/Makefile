WAIT := 200s
KUBECTL = kubectl --kubeconfig=.kubeconfig
ISTIOCTL = istioctl --kubeconfig=.kubeconfig
NAME = $$(echo $@ | cut -d "-" -f 2- | sed "s/%*$$//")
DOCKER = docker


install-istio-%: 
	$(ISTIOCTL) operator init 
	$(KUBECTL) apply -f default-profile.yaml

install-nginx-%: 
	$(KUBECTL) create ns ingress-nginx || true
	$(KUBECTL) label namespace ingress-nginx istio-injection=enabled || true
	$(KUBECTL) apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/cloud/deploy.yaml
	$(KUBECTL) wait --for=condition=Available deployment/ingress-nginx-controller -n ingress-nginx --timeout=160s

install-httpbin-%: install-istio-% 
	$(KUBECTL) patch svc ingress-nginx-controller -n ingress-nginx --patch "$$(cat patches/patch_ingresscontroller.yaml)"
	$(KUBECTL) patch svc istio-ingressgateway -n istio-system --patch "$$(cat patches/patch_ingressgateway.yaml)"
	$(KUBECTL) label namespace default istio-injection=enabled || true
	$(KUBECTL) apply -f httpbin.yaml
	$(KUBECTL) apply -f lua_filter.yaml

test-%:
	$(eval POD := $(shell $(KUBECTL) get pods -n istio-system -lapp=istio-ingressgateway -o=name))	
	$(KUBECTL) exec -ti $(POD) -n istio-system -- curl -v http://istio-ingressgateway/mycontext/get
	$(KUBECTL) exec -ti $(POD) -n istio-system -- curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJJc3N1ZXIiOiJJc3N1ZXIiLCJ1c3IiOiJKYXZhSW5Vc2UiLCJleHAiOjE2NzczMzI3ODIsImJyYW5jaCI6ImIxIiwiaWF0IjoxNjQ1Nzk2NzgyfQ.I3y5OuHxvW2VI46HCnfXENqvUhfAu11YyTKFUNvaXU8" http://istio-ingressgateway/mycontext/get 

delete-%:
	@kind delete cluster --name $(NAME)

env-%:
	@kind get kubeconfig --name $(NAME) > .kubeconfig

mp-env-%:
	cd multipass && ./kubeconfig.sh

clean:
	@kind get clusters | xargs -L1 -I% kind delete cluster --name %

list:
	@kind get clusters

create-%:
	-@sh/kind-registry.sh $(NAME) $(WAIT)

mp-create-%:
	cd multipass && ./install.sh

mp-delete-%:
	-@cd multipass
	-@reset.sh

.PHONY: mp-create-% create-% delete-% env-% clean list install-istio-% install-httpbin-%