
/*  PulseSensor Starter Project and Signal Tester
    The Best Way to Get Started  With, or See the Raw Signal of, your PulseSensor.com™ & Arduino.

    Here is a link to the tutorial
    https://pulsesensor.com/pages/code-and-guide

    WATCH ME (Tutorial Video):
    https://www.youtube.com/watch?v=RbB8NSRa5X4


  -------------------------------------------------------------
  1) This shows a live human Heartbeat Pulse.
  2) Live visualization in Arduino's Cool "Serial Plotter".
  3) Blink an LED on each Heartbeat.
  4) This is the direct Pulse Sensor's Signal.
  5) A great first-step in troubleshooting your circuit and connections.
  6) "Human-readable" code that is newbie friendly."

*/


//  Variables
int PulseSensorPurplePin = 0;        // Pulse Sensor PURPLE WIRE connected to ANALOG PIN 0
int LED13 = 13;   //  The on-board Arduion LED


int Signal;                // holds the incoming raw data. Signal value can range from 0-1024
int Threshold = 550;            // Determine which Signal to "count as a beat", and which to ingore.
int counter = 0;

// The SetUp Function:
void setup() {
  Serial.begin(9600);         // Set's up Serial Communication at certain speed.
}

// The Main Loop Function
void loop() {
  Signal = analogRead(PulseSensorPurplePin);  // Read the PulseSensor's value.

    Serial.println(Signal);                    // Send the Signal value to Serial Plotter.
  if (Signal > Threshold && counter == 0) {
    Serial.println("b");
    counter = 1;
  }
  if (counter == 1 && Signal < Threshold) {
    counter = 0;
  }
  delay(10);
}