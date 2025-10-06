(
Server.killAll;
s = Server.local();
s.reboot();
s.waitForBoot {
	"server booted".postln;
	s.freeAll;
	CmdPeriod.run;

	// Load all important files (maybe split htis up? We'll see)
	~projectDir = PathName(thisProcess.nowExecutingPath).parentPath;
	this.executeFile(~projectDir ++ "menu.scd");

	// Server is booted, we must initialize basic OSC server
	f = { |msg, time, addr, recvPort|
		// Handle message, play the correct song based on the message
		// received from the game
		"received!".postln;
		msg.postln;
		switch (msg[1])
		{'menu'} { ~playMusic.(); }
		{'stopall'} { ~stopMusic.(); }
	};

	o = OSCFunc(f, '/osc', nil, 9637);

	("Ready and listening on port " ++ o.recvPort).postln;

	// Send message indicating that the supercollider server is ready.
	// This message will be awaited in the loading screen.
	Routine({
		loop {
			1.0.wait;
		}
	}).play(AppClock);
};
)