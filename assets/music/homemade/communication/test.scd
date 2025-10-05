/*(
s = Server.local();
s.reboot();
s.waitForBoot {
	"Booted!".postln;

	SynthDef(\sin, {
		var snd, env;
		snd = SinOsc.ar(440.0);
		env = EnvGen.kr(Env([0, 1, 1, 0], [1, 1, 3], [1, 0, -1]),doneAction:2);
		snd = snd * env * -20.dbamp;
		Out.ar(0, snd);
	}).add;

	s.sync;

	1.do {
		Synth(\sin);
	};
}
)*/
// (
// s.plotTree;
// s.meter;
// )


// (
// Server.killAll;
// )

// lkmsdc
// This is very good - Now we just need to find out how to launch this from
// the game itself..
(
f = { |msg, time, addr, recvPort|
	1.do {
		Synth(\sin);
	};
}
)

(
o = OSCFunc(f, '/osc/address');
)

(
o.free;
)