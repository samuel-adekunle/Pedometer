// Here are some useful utility functions, which allow manipulation of the historic step data arrays.

export function getTotalSteps(step_histogram) {
	return Object.values(step_histogram).reduce((prev, curr) => prev + curr, 0);
}

export function getTodaysSteps(step_histogram) {
	let today_str = new Date().toISOString().slice(0, 10);
	let today_steps = step_histogram[today_str];

	return today_steps ? today_steps : 0;
}
