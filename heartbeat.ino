
/*  heartbeat.ino */


//  Variables
// Pulse Sensor PURPLE WIRE connected to ANALOG PIN 0
int pulseSensorPurplePin = 0;       
// signal holds the incoming raw data. signal value can range from 0-1024
int signal;                
// Threshold to use "count as a beat"
int Threshold = 550;      
// counter for hysteresis      
int counter = 0;

void setup() {
  Serial.begin(9600);       
}

void loop() {
  signal = analogRead(pulseSensorPurplePin);  
  Serial.println(signal);                    

  if (signal > Threshold && counter == 0) {
    Serial.println("b");
    counter = 1;
  }

  if (counter == 1 && signal < Threshold) {
    counter = 0;
  }

  delay(10);
}