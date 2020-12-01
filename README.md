    Apointment Project
    
    git clone https://github.com/kompir/app-pointment.git
    cd app-pointment/
    go mod tidy

    cd app-pointment/
    make

    1st bash
    cd app-pointment/notifier/
    node notifier.js

    2nd bash
    ./app-pointment/bin/server

    3rd bash
    ./app-pointment/bin/client create --title="Appoitment title" --message="Alert message." --duration=1m
    

 