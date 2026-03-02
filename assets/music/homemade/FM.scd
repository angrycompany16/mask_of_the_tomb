(
s.reboot;
s = Server.local;
)

(
s.plotTree;
s.meter;
s.scope;
)

s.scope;

// FM rule 1: As modulator amplitude increases, we get more audible sidebands, leading to a richer sound.
// FM rule 2: The distance between the sidebands is "equal" (proportional) to the modulator frequency.
// (FM rule 2.5:) The offset of the middle band is determined by the carrier frequency (but this is obvious.)
// FM rule 3: "Pitch clearness" tends to be greater when carrier frequency are both integer multiples of a "relatively large number". I.e. ratios like 2/1, 3/2, etc.

(
SynthDef(\fm, {
	arg freq=500, mRatio=1, cRatio=1, index=1, iScale=5, amp=0.2, atk=0.01, rel=3, pan=0, cAtk=4, cRel=(-4);
	var car, mod, env, iEnv, mod2;
	iEnv = EnvGen.kr(
		Env.new(
			[index, index * iScale, index],
			[atk, rel],
			[cAtk, cRel],
		)
	);

	env = EnvGen.kr(
		Env.perc(atk, rel, curve:[cAtk, cRel]), doneAction: Done.freeSelf
	);
	mod2 = SinOsc.ar(freq/5, mul:freq/10 * iEnv);
	mod = SinOsc.ar(freq*mRatio + mod2, mul:freq * mRatio * iEnv);
	car = SinOsc.ar((freq*cRatio) + mod) * env * amp;
	car = Pan2.ar(car, pan);
	Out.ar(0, car);
}).add;
)

(
Synth(\fm, [
	\freq, 60.midicps,
	\amp, 0.1,
	\mRatio, 5,
	\cRatio, 3,
	\index, 2,
	\iScale, 2,
	\atk, 0.01,
	\rel, 2.0,
]);
)
