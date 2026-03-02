(
s.reboot;
s = Server.local;
)

(
s.plotTree;
s.meter;
FreqScope.new;
)

(
s.scope;
)

// NOTE: All wavetables must be the same size! Otherwise it
// fails silently

(
~wt = Array.fill(4, {
	var numSegs = rrand(4, 20);
	Env(
		[0]++(({rrand(0, 1)}!(numSegs-1)) * [1, -1]).scramble++[0],
		{exprand(1, 20)}!numSegs,
		{rrand(-20, 20)}!numSegs,
	).asSignal(1024).asWavetable;
});
)

~buf = Buffer.allocConsecutive(4, s, 2048);

(
~buf.do({ |buf, i|
	buf.loadCollection(~wt[i]);
});
)

(
SynthDef.new(\vosc, {
	arg out=0, buf=0, numBufs=2, freq=300, amp=(-20);
	var sig, bufpos, detuneSig;
	detuneSig = LFNoise1.kr(0.2!8).bipolar(0.2).midiratio;
	bufpos = buf + LFNoise1.kr(0.5).range(0, numBufs-1);
	sig = VOsc.ar(bufpos, freq*detuneSig);
	sig = Splay.ar(sig);
	sig = LeakDC.ar(sig);
	sig = sig * amp.dbamp;
	Out.ar(out, sig);
}).add;
)

Synth(\vosc, [\buf, ~buf[0], \freq, 24.midicps, \numBufs, 4, \amp, -20]);


s.freeAll;