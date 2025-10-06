(
this.executeFile(~projectDir ++ "pads/final.scd");

~ping = { | frequencies |
	Routine {
		loop {
			rrand(5.0, 8.0).wait;
			11.do {
				if (frequencies.notNil) {
					Synth(\ping,  [freq: frequencies.choose * 2, amp: -28.dbamp, pan: rrand(-1.0, 1.0)]);
				};
				exprand(0.01, 0.5).wait;
			};
		};
	};
};

~fm = { | frequencies |
	Routine {
		loop {
			rrand(5.0, 8.0).wait;
			6.do {
				if (frequencies.notNil) {
					Synth(\fm,  [freq: frequencies.choose, amp: -28.dbamp, pan: rrand(-1.0, 1.0)]);
				};
				exprand(0.1, 0.5).wait;
			};
		};
	};
};

~pad = { | frequency |
	Routine {
		var synth;
		synth = Synth(\pad, [freq: frequency, amp: -20.dbamp] );
		rrand(1.0, 4.0).wait;
		synth.set(\gate, 0);
	};
};

~playMusic = {
	var frequencies, chords;

	~stopMusic.();

	chords = [
		[3, 7, 12, 14, 15, 24, 26],
		[3, 7, 12, 14, 15, 24, 26],
		[3, 7, 12, 14, 15, 24, 26, 27, 29],
		[3, 7, 12, 14, 15, 24, 26, 27, 29],
		[0, 15, 24, 26, 27, 31],
	];

	Synth.tail(nil, \fx);

	~pingRoutine = ~ping.().play;
	~fmRoutine = ~fm.().play;

	~padRoutine = Routine {
		loop {
			chords.do { |chord|
				frequencies = (39 + chord +.x [0, 12]).midicps;

				frequencies.scramble.do { |frequency|
					~pad.(frequency).play;
					rrand(0.3, 1.0).wait;
				};
				rrand(5.0, 10.0).wait;
			};
		};
	}.play;
};

~stopMusic = {
	if (~pingRoutine.notNil) { ~pingRoutine.stop; ~pingRoutine = nil; };
	if (~fmRoutine.notNil) { ~fmRoutine.stop; ~fmRoutine= nil; };
	if (~padRoutine.notNil) { ~padRoutine.stop; ~padRoutine = nil; };
};
)