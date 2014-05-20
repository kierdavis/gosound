package clfourier

var CoreSource = `
__kernel void compute_cell(
	__global float *input,
	__global float *freqs,
	__global float *reals,
	__global float *imags,
	const unsigned int log_input_size,
	const unsigned int count,
) {
	unsigned int i = get_global_id(0);
	if (i < count) {
		unsigned int freq_idx = i >> log_input_size;
		unsigned int sample_idx = i & ((1 << log_input_size) - 1);
		float sample = input[sample_idx];
		float phase = -2 * M_PI * freqs[freq_idx] * ((float) sample_idx);
		float cos_phase, sin_phase = sincos(phase, &cos_phase);
		reals[i] = sample * cos_phase;
		imags[i] = sample * sin_phase;
	}
}

__kernel void partial_sum(
	__global float *reals,
	__global float *imags,
	const unsigned int round,
	count unsigned int count,
) {
	unsigned int i = get_global_id(0);
	if (i < count) {
		unsigned int idx1 = i << round;
		unsigned int idx2 = idx1 + (1 << (round - 1));
		reals[idx1] += reals[idx2];
		imags[idx1] += imags[idx2];
	}
}

__kernel void average(
	__global float *reals,
	__global float *imags,
	__global float *output,
	const unsigned int log_input_size,
	const unsigned int count,
) {
	unsigned int i = get_global_id(0);
	if (i < count) {
		unsigned int idx = i << log_input_size;
		float re = reals[idx];
		float im = imags[idx];
		output[i] = sqrt(re*re + im*im) / ((float) (1 << log_input_size));
	}
}
`
