import { interpolatePuBu } from "d3-scale-chromatic";
import Head from "next/head";
import { useRouter } from "next/router";
import { useState } from "react";
import { Button, ProgressBar, Table } from "react-bootstrap";
import { Pie } from "react-chartjs-2";
import Confetti from "react-confetti";
import useSWR, { mutate } from "swr";
import interpolateColors from "../utils/color-selection";
import { API_KEY, BASE_URL } from "../utils/fetcher";
import SideNavPanel from "./SideNavPanel";

function PieChart({ steps, target }) {
	let users = Object.keys(steps);
	let steps_data = Object.values(steps);
	const data_sum = steps_data.reduce((a, b) => {
		return a + b;
	});

	let incomplete = target - data_sum;
	if (incomplete > 0) {
		users.push("Incomplete");
		steps_data.push(target - data_sum);
	}
	const colorScale = interpolatePuBu;
	const colorRangeInfo = {
		colorStart: 0,
		colorEnd: 1,
		useEndAsStart: true,
	};
	let colors = interpolateColors(users.length, colorScale, colorRangeInfo);
	const data = {
		labels: users,
		datasets: [
			{
				backgroundColor: colors,
				hoverBackgroundColor: colors,
				data: steps_data,
			},
		],
	};
	const options = {
		responsive: true,
		legend: {
			display: true,
		},
		maintainAspectRatio: true
	};
	return <Pie data={data} options={options} />;
}

function ChallengesItem({
	challenge_name,
	current,
	target,
	steps,
}) {
	const router = useRouter();

	const progress = Math.min(1, current / target) * 100;

	const handleDelete = (e) => {
		e.preventDefault();
		fetch(`${BASE_URL}/challenges/${challenge_name}`, {
			method: "DELETE",
			headers: {
				"x-api-key": API_KEY,
			},
		}).then(() => {
			mutate(`${BASE_URL}/challenges`);
			router.replace("/");
		});
	};

	return (
		<>
			<div style={{ display: "flex", alignItems: "center" }}>
				<h5 style={{ flex: 2, marginBottom: 0 }}>Challenge Progress: </h5>
				<ProgressBar style={{ flex: 8, marginTop: 0, marginRight: "1%" }} now={progress} />
				<div style={{ flex: 1 }}>
					{current}/{target}
				</div>
			</div>
			<br />
			<div>
				<PieChart steps={steps} target={target} />
			</div>
			<h4>Steps:</h4>
			<br />
			<div>
				<Table bordered>
					<thead>
						<tr>
							<th>User Name</th>
							<th>Number of Steps</th>
						</tr>
					</thead>
					<tbody>
						{Object.keys(steps).map((userName) => <tr key={userName}>
							<td>{userName}</td>
							<td>{steps[userName]}</td>
						</tr>)}
					</tbody>
				</Table>
			</div>
			<br />
			{challenge_name !== "default" && (
				<div>
					<Button variant="danger" onClick={handleDelete}>Delete challenge</Button>
				</div>
			)}
			<br />
		</>
	);
}

function ChallengesPage({ initialData, challengeName }) {
	const { data } = useSWR(`${BASE_URL}/challenges/${challengeName}`, {
		initialData,
	});

	const [completedChallenge, setCompletedChallenge] = useState(false);

	if (data && data.data) {
		const {
			challenge_name,
			current,
			target,
			start_date,
			end_date,
			steps,
		} = data.data;

		if (!completedChallenge) {
			if (current >= target) {
				setCompletedChallenge(true);
			}
		}

		return <>
			<Head>
				<title>{challengeName} - Pedometer</title>
			</Head>
			{completedChallenge && <Confetti />}
			<div>
				<SideNavPanel active_item="challenges" />
				<div className="page-content">
					<h1 className="page-title">{challenge_name}</h1>
					<br />
					<ChallengesItem
						challenge_name={challenge_name}
						current={current}
						target={target}
						start_date={start_date}
						end_date={end_date}
						steps={steps}
					/>
				</div>
			</div>
		</>;
	}
}

export default ChallengesPage;
