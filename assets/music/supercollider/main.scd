Server.killAll;
s = Server.local();
s.reboot();
s.waitForBoot {
	s.freeAll;
	CmdPeriod.run;

	// Load all important files (maybe split htis up? We'll see)
	~projectDir = PathName(thisProcess.nowExecutingPath).parentPath;
	this.executeFile(~projectDir ++ "menu.scd");

	// Server is booted, we must initialize basic OSC server
	f = { |msg, time, addr, recvPort|
		// Handle message, play the correct song based on the message
		// received from the game
		msg.postln;
		if (msg[1] == 'menu') {
			~playMusic.value.play();
		};
	};

	o = OSCFunc(f, '/osc');

	// Send message indicating that the supercollider server is ready.
	// This message will be awaited in the loading screen.
	loop {
		1.0.wait;
	};

	o.free;
}