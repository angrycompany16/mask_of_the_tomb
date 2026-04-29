(
s = Server.local;
s.reboot;
)

s.scope;
s.plotTree;

~b = Buffer.read(s, "/home/angrycompany16/gamedev/go/mask_of_the_tomb/assets/music/homemade/samples/beat.wav");
Buffer.freeAll;
~b.play;

(
SynthDef(\fm, {
	arg freq=500, mRatio=1, cRatio=1, index=1, iScale=5, cAtk=4, cRel=(-4), amp=0.2,atk=0.01, rel=3, pan=0, out=0;
	var car, mod, env, iEnv;
	iEnv = EnvGen.kr(
		Env.new(
			[index, index * iScale, index],
			[atk, rel],
			[cAtk, cRel],
		)
	);
	env = EnvGen.kr(Env.perc(atk, rel, curve:[cAtk, cRel]), doneAction:2);
	mod = SinOsc.ar(freq * mRatio, mul:freq * mRatio * iEnv);
	car = SinOsc.ar(freq * cRatio + mod) * env * amp;
	car = Pan2.ar(car, pan);
	Out.ar(out, car);
}).add;

// 2-channel reverb
SynthDef(\reverb, {
	arg in, predelay=0.1, revtime=1.8, lpf=4500, mix=0.15, amp=1, out=0;
	var dry, wet, temp, sig;
	dry = In.ar(in, 2);
	temp = In.ar(in, 2);
	wet = 0;
	// predelay essentially becomes like room size
	temp = DelayN.ar(temp, 0.2, predelay);
	16.do {
		temp = AllpassN.ar(temp, 0.05, {Rand(0.001, 0.05)}!2, revtime);
		temp = LPF.ar(temp, lpf);
		wet = wet + temp;
	};
	sig = XFade2.ar(dry, wet, mix*2-1, amp);
	Out.ar(out, sig);
}).add;

SynthDef(\grainBuf, {
	arg bufnum, density=1, trig=0, grainDur=1, rate=1, pos=0, pan=0, lpf=3000, amp=1;
	var sig, trigger;

	trigger = Select.ar(trig, [Dust.ar(density), Impulse.ar(density)]);

	sig = GrainBuf.ar(
		2,
		trigger,
		grainDur,
		bufnum,
		rate,
		pos,
		2,
		pan,
		-1,
		512,
	);

	sig = LPF.ar(sig, lpf);
	sig = sig * amp;

	Out.ar(0, sig);
}).add;
)

~reverbBus = Bus.audio(s, 2);

(
~reverbSynth = Synth(\reverb, [
	\in, ~reverbBus,
	\predelay, 0.1,
	\revtime, 1.0,
	\lpf, 4500,
	\mix, 0.3,
	\amp, 1,
]);
)

(
bpm = 100;
~clock = TempoClock.default();
~clock.tempo = bpm / 60;
)

(
var drumSpeed = 0.25;
var grainDur = 4;

// Bruh, should have just done this with a regular buffer synth
// tbh
Pdef(\drums,
	Pbind(
		\instrument, \grainBuf,
		\dur, inf,

		\bufnum, ~b,
		\density, ~clock.tempo * drumSpeed / grainDur,
		\trig, 1,
		\pos, 0.7,
		\grainDur, grainDur / drumSpeed,
		\rate, drumSpeed,
		\lpf, 3500,
		\amp, 0.6,
	)
);
)

Pbindef(\drums).play;
Pbindef(\drums).stop;

(
Synth(\fm, [
	\freq, 60.midicps,
	\mRatio, 1,
	\cRatio, 1,
	\index, 2,
	\iScale, 2,
	\rel, 4.0,
	\cAtk, 4,
	\cRel, -4,
	\amp, 0.1,
	\out, ~reverbBus,
]);
)

(
Synth(\grainBuf, [
	\bufnum, ~b,
	\density, 1/4,
	\trig, 1,
	\pos, 0.32,
	\dur, 3.42 * 4,
	\rate, 1/4,
	\lpf, 1000,
	\amp, 0.6,
]);
)














