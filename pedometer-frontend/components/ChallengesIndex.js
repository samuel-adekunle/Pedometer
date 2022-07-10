import { Button, SelectDate, TextField } from "@duik/it";
import Head from "next/head";
import Link from "next/link";
import { useState } from "react";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import ProgressBar from "react-bootstrap/ProgressBar";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import Tooltip from "react-bootstrap/Tooltip";
import { ImCancelCircle } from "react-icons/im";
import { MdAdd } from "react-icons/md";
import useSWR, { mutate } from "swr";
import { API_KEY, BASE_URL } from "../utils/fetcher";
import SideNavPanel from "./SideNavPanel";


function ChallengesForm() {
	const [challengeName, setChallengeName] = useState("");
	const [startDate, setStartDate] = useState(new Date());
	const [endDate, setEndDate] = useState(new Date());
	const [target, setTarget] = useState(1);
	const [expand_form, set_expand_form] = useState(false);

	const handleSubmit = (e) => {
		e.preventDefault();

		fetch(`${BASE_URL}/challenges/${challengeName}`, {
			method: "POST",
			headers: {
				"x-api-key": API_KEY,
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				challenge_name: challengeName,
				start_date: new Date(startDate.setDate(startDate.getDate() + 1)).toISOString().split("T")[0],
				end_date: new Date(endDate.setDate(endDate.getDate() + 1)).toISOString().split("T")[0],
				target: Number(target)
			})
		}).then(() => {
			mutate(`${BASE_URL}/challenges`);
			set_expand_form(false);
			setChallengeName("");
			setEndDate(new Date());
			setStartDate(new Date);
			setTarget("");
		});
	};

	function toggleClick() {
		set_expand_form(!expand_form);
	}

	if (!expand_form) {
		return (
			<tr className="new-user-row clickable">
				<td colSpan="3" onClick={toggleClick}><MdAdd /> Create New Challenge</td>
			</tr>
		);
	} else {
		return (<>
			<tr className="new-user-row clickable">
				<td style={{ color: "rgb(246, 96, 22)" }} colSpan="3" onClick={toggleClick}><ImCancelCircle /> Cancel</td>
			</tr>
			<tr>
				<td colSpan="3">
					<Container>
						<Row>
							<Col>
								<TextField value={challengeName} onChange={(e) => setChallengeName(e.target.value)} label="Challenge Name" placeholder="e.g. my-challenge" />
							</Col>
							<Col>
								<TextField value={target} onChange={(e) => setTarget(e.target.value)} min={1} type="number" label="Step Goal" placeholder="e.g. 1000" />
							</Col>
							<Col>
								<label>Start Date</label><br />
								<SelectDate value={startDate} onDateChange={setStartDate} minDate={new Date()} />
							</Col>
							<Col>
								<label>End Date</label><br />
								<SelectDate value={endDate} onDateChange={setEndDate} minDate={new Date()} />
							</Col>
							<Col>
								<div className="add-user-button-wrapper">
									<Button primary block onClick={handleSubmit}>
										Add Challenge +
                  </Button>
								</div>
							</Col>
						</Row>
					</Container>
				</td>
			</tr>
		</>);
	}

}

function ChallengesIndex({ initialData }) {
	const { data } = useSWR(`${BASE_URL}/challenges`, { initialData });

	const progress = ({ current, target }) => (current / target);

	if (data) {
		// Sort challenges by progress
		var sorted_data = data.data.sort((a, b) => {
			return Math.min(1, progress(b) - progress(a));
		});

		return <>
			<Head>
				<title>Challenges - Pedometer</title>
			</Head>
			<SideNavPanel active_item="challenges" />
			<div className="page-content">
				<h1 className="page-title">Challenges</h1>
				<p>
					From this page, the challenges can be browsed. Just click on the challenge in the list to view detailed information.
        </p>
				<br />
				<Table bordered>
					<thead>
						<tr>
							<th>Challenge</th>
							<th>Step Target</th>
							<th>Progress</th>
						</tr>
					</thead>
					<tbody>
						{
							sorted_data.map((challenge) => (
								<Link key={challenge.challenge_name} href={`/challenges/${challenge.challenge_name}`}>
									<tr className="clickable">
										<td>{challenge.challenge_name}</td>
										<td>{challenge.target}</td>
										<td colSpan="2">
											<OverlayTrigger
												key="top"
												placement="top"
												overlay={
													<Tooltip id="tooltip-top">
														Deadline on {challenge.end_date}
													</Tooltip>
												}>
												<ProgressBar variant="primary" now={progress(challenge) * 100} />
											</OverlayTrigger>
										</td>
									</tr>
								</Link>
							))
						}
						<ChallengesForm />
					</tbody>
				</Table>
			</div>
		</>;
	}

	return <h1>Loading...</h1>;
}

export default ChallengesIndex;
