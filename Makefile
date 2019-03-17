bin:
	go build -o terraform-provider-discord

clean:
	rm -rf ./.terraform

test-integration: bin clean
	TF_LOG=DEBUG terraform init ./test && \
	TF_LOG=DEBUG terraform plan ./test
