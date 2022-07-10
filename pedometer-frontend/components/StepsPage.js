import TextField from "@duik/text-field";
import Head from "next/head";
import { useRouter } from "next/router";
import React, { useState } from "react";
import { Button, Card, ProgressBar, Table } from "react-bootstrap";
import { Bar } from "react-chartjs-2";
import { ImCancelCircle } from "react-icons/im";
import { MdAdd } from "react-icons/md";
import useSWR, { mutate } from "swr";
import { API_KEY, BASE_URL } from "../utils/fetcher";
import { getTotalSteps } from "../utils/historic_data_manipulation";
import SideNavPanel from "./SideNavPanel";

function BarChart({ count }) {
	let dates = Object.keys(count).reverse();
	let count_data = Object.values(count).reverse();
	const data = {
		labels: dates,
		datasets: [
			{
				label: "Number of steps",
				backgroundColor: "rgba(72,159,181,0.2)",
				borderColor: "rgba(72,159,181,1)",
				borderWidth: 1,
				hoverBackgroundColor: "rgba(72,159,181,0.4)",
				hoverBorderColor: "rgba(72,159,181,1)",
				data: count_data,
			},
		],
	};
	return (
		<Bar
			data={data}
			width={100}
			height={50}
			options={{
				maintainAspectRatio: true,
			}}
		/>
	);
}

// NOTE: function to change daily target
//   const handleSubmit = (e) => {
//     e.preventDefault();

//     fetch(`${BASE_URL}/steps/${userName}`, {
//       method: "POST",
//       headers: {
//         "x-api-key": API_KEY,
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({
//         user_name: userName,
//         daily_target: Number(dailyTarget),
//       }),
//     }).then(() => {
//       mutate(`${BASE_URL}/steps/${userName}`);
//       setDailyTarget(0);
//     });
//   };

function StepsChallengeItem({ challengeName, userName }) {
	const router = useRouter();

	const handleRemoveChallenge = (e) => {
		e.preventDefault();

		fetch(`${BASE_URL}/steps/${userName}/${challengeName}`, {
			method: "DELETE",
			headers: {
				"x-api-key": API_KEY,
			},
		}).then(() => {
			mutate(`${BASE_URL}/steps/${userName}/${challengeName}`);
			mutate(`${BASE_URL}/steps/${userName}`);
		});
	};

	return <>
		<tr className="clickable" onClick={() => router.push(`/challenges/${challengeName}`)}>
			<td>
				<div style={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center"
				}}>
					{challengeName}
					{challengeName !== "default" && <Button variant="danger" onClick={handleRemoveChallenge}>Remove</Button>}
				</div>
			</td>
		</tr>
	</>;
}

function StepsItem({ user_name, challenges, count, daily_target }) {
	const router = useRouter();
	const [expandForm, setExpandForm] = useState(false);

	const [newChallengeName, setNewChallengeName] = useState("");

	const today = count[new Date().toISOString().slice(0, 10)];
	const progress = Math.min(1, (today / daily_target)) * 100;
	const averageSteps = getTotalSteps(count) / Object.keys(count).length;
	const averageCalories = averageSteps * 0.04;

	const handleDelete = (e) => {
		e.preventDefault();
		fetch(`${BASE_URL}/steps/${user_name}`, {
			method: "DELETE",
			headers: {
				"x-api-key": API_KEY,
			},
		}).then(() => {
			mutate(`${BASE_URL}/steps`);
			router.replace("/");
		});
	};

	const handleAdd = (e) => {
		e.preventDefault();

		fetch(`${BASE_URL}/steps/${user_name}/${newChallengeName}`, {
			method: "POST",
			headers: {
				"x-api-key": API_KEY,
			},
		}).then(() => {
			mutate(`${BASE_URL}/steps/${user_name}`);
			mutate(`${BASE_URL}/challenges/${newChallengeName}`);
		});
	};

	return (
		<>
			<div style={{ display: "flex", alignItems: "center" }}>
				<h5 style={{ flex: 1, marginBottom: 0 }}>Daily Target: </h5>
				<ProgressBar style={{ flex: 8, marginTop: 0, marginRight: "1%" }} now={progress} />
				<div style={{ flex: 1 }}>
					{today}/{daily_target}
				</div>
			</div>
			<br />
			<div style={{ display: "flex", justifyContent: "space-between" }}>
				<Card style={{ width: 210 }}>
					<Card.Body>
						<div>Daily Target</div>
						<div>{daily_target}</div>
					</Card.Body>
				</Card>
				<Card style={{ width: 210 }}>
					<Card.Body>
						<div>Max Daily Steps</div>
						<div>{Math.max(...Object.values(count))}</div>
					</Card.Body>
				</Card>
				<Card style={{ width: 210 }}>
					<Card.Body>
						<div>Min Daily Steps</div>
						<div>{Math.min(...Object.values(count))}</div>
					</Card.Body>
				</Card>
				<Card style={{ width: 210 }}>
					<Card.Body>
						<div>Avg Daily Steps</div>
						<div>{Math.round(averageSteps)}</div>
					</Card.Body>
				</Card>
				<Card style={{ width: 210 }}>
					<Card.Body>
						<div>Avg Calories Burnt</div>
						<div>{Math.round(averageCalories)} kcal</div>
					</Card.Body>
				</Card>
			</div>
			<br />
			<div>
				<BarChart count={count} />
			</div>
			<br />
			<h4>
				Challenges:
      </h4>
			<div>
				<Table bordered>
					<thead>
						<tr>
							<th>
								Challenge Name
              </th>
						</tr>
					</thead>
					<tbody>
						{challenges.map((challengeName) => <React.Fragment key={challengeName}>
							<StepsChallengeItem
								challengeName={challengeName}
								userName={user_name}
							/>
						</React.Fragment>)}
						<tr className="new-user-row clickable">
							{!expandForm
								? <td onClick={() => setExpandForm((prev) => !prev)}>
									<MdAdd /> Enrol in new challenge
                </td>
								: <td style={{ color: "rgb(246, 96, 22)" }} onClick={() => setExpandForm((prev) => !prev)}>
									<ImCancelCircle /> Cancel
                </td>
							}
						</tr>
						{expandForm && <tr>
							<td>
								<div style={{ display: "flex", justifyContent: "center" }}>
									<div style={{ paddingRight: "1%", width: "50%" }}>
										<TextField
											label="Challenge Name"
											placeholder="my-challenge"
											value={newChallengeName}
											onChange={(e) => setNewChallengeName(e.target.value)} />
									</div>
									<div className="add-user-button-wrapper">
										<Button variant="primary" onClick={handleAdd}>Enrol in challenge <MdAdd /></Button>
									</div>
								</div>
							</td>
						</tr>}
					</tbody>
				</Table>
			</div>
			<br />
			<div>
				<Button variant="danger" onClick={handleDelete}>Delete user</Button>
			</div>
			<br />
		</>
	);
}

function StepsPage({ initialData, userName }) {
	const { data } = useSWR(`${BASE_URL}/steps/${userName}`, { initialData });

	if (data && data.data) {
		const { user_name, daily_target, challenges, count } = data.data;
		return (
			<>
				<Head>
					<title>{userName} - Pedometer</title>
				</Head>
				<SideNavPanel active_item="" />
				<div className="page-content">
					<h1 className="page-title">
						{user_name}
					</h1>
					<br />
					<StepsItem
						user_name={user_name}
						daily_target={daily_target}
						challenges={challenges}
						count={count}
					/>
				</div>
			</>
		);
	} else {
		return <h2>Loading...</h2>;
	}
}

export default StepsPage;
