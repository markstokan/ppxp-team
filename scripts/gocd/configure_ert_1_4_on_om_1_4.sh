source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake ert:configure[$RELENG_ENV,1.4,1.4]'
