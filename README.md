# heartbpm

Its 2020 and entirely possible to [control synthesizers with your thoughts](https://www.youtube.com/watch?v=qSKBtEBRWi4), so of course its entirely possible to control your synthesizers with your heart pulse. In fact, its entirely possible to control your synthesizers with your heart rate with only ~$50 of equipment. This repository contains the code and instructions to hook your blood to your music.


This system is really simple and will work with basically any modern synthesizer (as long as it has MIDI). Here's a little schematic:

<p align="center">
<img src="https://schollz.com/img/schematic1.PNG" alt="A basic schematic of the setup." width=550px>
</p>

A finger is connected to a pulse sensor. The pulse sensor is connected to an Arduino. The Arduino can detect a pulse and send serial data to a computer. The computer runs a server that reads the serial data which updates a web page. The web page uses Web MIDI to send out the MIDI clock to all connected synths! *Music* ensues.

Here's a demo of me playing with three instruments sequenced to the tempo from my heart rate:


<p align="center"><a href="https://www.youtube.com/watch?v=o_GRabvE5Pw"><img src="https://img.youtube.com/vi/o_GRabvE5Pw/0.jpg" alt="Demo of playing" style="max-width:200px;"></a></p>

The following instructions will give you a more detailed instruction.

## Instructions

### Stuff to buy.

You'll be able to get started with $50 of equipment, plus computer and MIDI instrument.

<ul>
	<li>
		<a target="_blank" href="https://www.amazon.com/gp/product/B008GRTSV6/ref=as_li_tl?ie=UTF8&camp=1789&creative=9325&creativeASIN=B008GRTSV6&linkCode=as2&tag=scholl-20&linkId=273d395910f0ec3a1c1da85d779214fb">
			Arduino ($23)
		</a>
	</li>
	<li>
		<a target="_blank" href="https://www.amazon.com/gp/product/B01CPP4QM0/ref=as_li_tl?ie=UTF8&camp=1789&creative=9325&creativeASIN=B01CPP4QM0&linkCode=as2&tag=scholl-20&linkId=e3026cfdfee3a180a780748bd1dce2e4">
			Pulse sensor ($25)
		</a>
	</li>
	<li>
		<a target="_blank" href="https://www.amazon.com/gp/product/B06XQYN77L/ref=as_li_tl?ie=UTF8&camp=1789&creative=9325&creativeASIN=B06XQYN77L&linkCode=as2&tag=scholl-20&linkId=ef5f1fd6a1f83b3fbbcabd242e4ff4cc">
			Audio isolator (optional, $9)
		</a>
	</li>
	<li>
		Synthesizer with MIDI (basically *any* synthesizer made in the last 10 years)
	</li>
	<li>Computer</li>
</ul>

### Get the code.

You can [clone the code from Github](https://github.com/schollz/heartbpm) or you can [download it directly](https://github.com/schollz/heartbpm/archive/master.zip).

### Setup the hardware.

Simply attach a pulse sensor to the Arduino, plugging the data channel into `A0`.

<p align="center">
<img alt="Connecting the pulse sensor to the Arduino" src="https://schollz.com/img/heartbeat_bb_POoEsoKAle.jpg" width=550px/>
</p>

Now use the [Arduino IDE](https://www.arduino.cc/en/main/software) to upload the [`heartbeat.ino`](https://github.com/schollz/heartbpm/blob/master/heartbeat.ino). Once that is uploaded, you can keep the Arduino connected and move on to the software.

### Setup the software.

First make sure you have Go installed on your computer You can [download Go here](https://golang.org/dl/). Now you can go into the `heartbpm` code and simply run in a terminal:

```bash
$ go build -v
```

Now you'll have an executable `heartbpm` in that directory. You can simply run it with

```
$ ./heartbpm --com COMPORT
[info]  2020/06/13 07:48:06 listening on :8054
```

Make sure you find your `COMPORT`. The easiest way to find your com port is to look at it from the Arduino IDE.


<p align="center">
<img alt="Finding the name of the COM port." src="https://schollz.com/img/comport.jpg" width=550px/>
</p>

Now, connect the pulse sensor to your finger and open a *Chrome* web browser to http://localhost:8054. You should be able to see some data coming out. Note: this only works with Chrome and Opera web browsers since they have Web MIDI implementations.

### Just add synthesizers.

To get it working with synthesizers, just attach your synthesizer via MIDI to your computer, Most modern synthesizers have USB which doubles as a MIDI connection, so just connect by USB! 

Connect your finger to the pulse sensor and then reload the web page. The web page will *automatically* detect the synthesizers and you can see which ones are attached. 


<p align="center">
<img alt="Screenshot of gathering data with NTS-1 attached." src="https://schollz.com/img/heratbpmscreen.png" width=550px/>
</p>

Your average heart rate will be detected and sent to the electronic instruments as the clock. The website is getting the averaged pulse data calculated BPM from the server. This BPM is used to send MIDI clock signals every 1/24th of a quarter note, which is the standard for setting tempos on devices.



### Useful notes.

- If you are having trouble getting your pulse connected, try moving the pulse sensor just a little bit along your finger (1/8"). Sometimes if its not over your vein it will be harder to detect. Once its in the right spot it will stay there nicely, though!
- If you are getting background noise its because the synthesizers are powered off your computer which (unless is battery operated) tends to have a noisy power supply. I highly recommend getting an [audio isolator](https://www.amazon.com/gp/product/B06XQYN77L/ref=as_li_tl?ie=UTF8&camp=1789&creative=9325&creativeASIN=B06XQYN77L&linkCode=as2&tag=scholl-20&linkId=ef5f1fd6a1f83b3fbbcabd242e4ff4cc) to remove this noise.

### Enjoy!

Hope this is useful for you, and hope you can create something new! If you are interested in the music I've created, check out [my Bandcamp](https://infinitedigits.bandcamp.com) or just search `infinite digits` on any streaming platform.

# License

MIT