gcloud projects get-iam-policy 1-tg  \
--flatten="bindings[].members" \
--format='table(bindings.role)' \
--filter="bindings.members:service-kozub-tg@containerregistry.iam.gserviceaccount.com"