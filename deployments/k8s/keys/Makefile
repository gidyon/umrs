ifdef version
	VERSION=-$(version)
else
	VERSION=
endif

create_account_tls:
	kubectl create secret generic account-tls$(VERSION) --from-file=cert=accounts.umrs.net/cert.pem --from-file=key=accounts.umrs.net/key.pem && echo $(VERSION)
create_ledger_tls:
	kubectl create secret generic ledger-tls$(VERSION) --from-file=cert=ledger.umrs.net/cert.pem --from-file=key=ledger.umrs.net/key.pem
create_employment_tls:
	kubectl create secret generic employment-tls$(VERSION) --from-file=cert=employments.umrs.net/cert.pem --from-file=key=employments.umrs.net/key.pem
create_hospital_tls:
	kubectl create secret generic hospital-tls$(VERSION) --from-file=cert=hospitals.umrs.net/cert.pem --from-file=key=hospitals.umrs.net/key.pem
create_insurance_tls:
	kubectl create secret generic insurance-tls$(VERSION) --from-file=cert=insurances.umrs.net/cert.pem --from-file=key=insurances.umrs.net/key.pem
create_umrs_tls:
	kubectl create secret generic umrs-tls$(VERSION) --from-file=cert=umrs.net/cert.pem --from-file=key=umrs.net/key.pem
create_notification_tls:
	kubectl create secret generic notification-tls$(VERSION) --from-file=cert=notifications.umrs.net/cert.pem --from-file=key=notifications.umrs.net/key.pem
create_patient_tls:
	kubectl create secret generic patient-tls$(VERSION) --from-file=cert=patients.umrs.net/cert.pem --from-file=key=patients.umrs.net/key.pem
create_permission_tls:
	kubectl create secret generic permission-tls$(VERSION) --from-file=cert=permissions.umrs.net/cert.pem --from-file=key=permissions.umrs.net/key.pem
create_treatment_tls:
	kubectl create secret generic treatment-tls$(VERSION) --from-file=cert=treatments.umrs.net/cert.pem --from-file=key=treatments.umrs.net/key.pem

create_ledger_enc_key:
	kubectl create secret generic ledger-secretkey$(VERSION) --from-literal=secret-key='J+JuQg7KVAAN87w/ybnFTv5/4417BjnodZPy7pyMO40='

init_tls_secret: create_treatment_tls create_permission_tls create_patient_tls create_notification_tls create_umrs_tls create_insurance_tls create_hospital_tls create_employment_tls create_ledger_tls create_account_tls create_ledger_enc_key