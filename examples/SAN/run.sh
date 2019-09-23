go run main.go -platform=input/virt_setup.xml \
		-packet=input/new_format.json \
		-controlling_mode=0 \
		-sim_run=1 \
		-atm_dep=input/temp.json \
		-atm_control=input/atm_control.json \
		-file_amount_w=1 \
		-file_size_w=100MB..100MB \
		-output=output.json \
		-num_jobs_config=input/num_jobs.json \
        	-num_jobs=1