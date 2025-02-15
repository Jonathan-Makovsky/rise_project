@echo off

:: Step 1: Set the path to the docker-compose.yml file (change the path if necessary)
set COMPOSE_PATH=C:\Users\Yonatan\Desktop\Rise\rise_project\setup

:: Step 2: Start the Docker containers using Docker Compose
echo Starting Docker containers...
docker-compose -f %COMPOSE_PATH%\docker-compose.yml up -d --build

:: Step 3: Wait for the services to be up and running (you can adjust the sleep time)
echo Waiting for the services to be ready...
timeout /t 10 /nobreak >nul

:: Step 4: Run the test file for the backend API endpoints (adjust the path to the test file)
echo Running the tests...
go test -v . -count=1

:: Optionally, you can check the exit code of the tests to decide what to do next
if %ERRORLEVEL% equ 0 (
  echo Tests passed successfully!
) else (
  echo Tests failed. Check the output for errors.
)

pause
