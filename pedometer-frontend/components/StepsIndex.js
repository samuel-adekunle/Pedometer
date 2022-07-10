import { Button, TextField } from "@duik/it";
import Head from "next/head";
import Link from "next/link";
import { useState } from "react";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import { ImCancelCircle } from "react-icons/im";
import { MdAdd } from "react-icons/md";
import useSWR, { mutate } from "swr";
import { API_KEY, BASE_URL } from "../utils/fetcher";
import { getTodaysSteps } from "../utils/historic_data_manipulation";
import SideNavPanel from "./SideNavPanel";

function StepsForm() {
	const [userName, setUserName] = useState("");
	const [count, setCount] = useState(0);
	const [expand_form, set_expand_form] = useState(false);

	const handleSubmit = (e) => {
		e.preventDefault();

		fetch(`${BASE_URL}/steps/${userName}`, {
			method: "POST",
			headers: {
				"x-api-key": API_KEY,
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				user_name: userName,
				count: {
					[new Date().toISOString().slice(0, 10)]: Number(count)
				}
			})
		}).then(() => {
			mutate(`${BASE_URL}/steps`);
			setUserName("");
			setCount(0);
			set_expand_form(false);
		});
	};

	function toggleClick() {
		set_expand_form(!expand_form);
	}

	if (!expand_form) {
		return (
			<tr className="new-user-row clickable">
				<td colSpan="3" onClick={toggleClick}><MdAdd /> Create New User</td>
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
							<Col><TextField label="Username" placeholder="Username" value={userName} onChange={(e) => setUserName(e.target.value)} />
							</Col>
							<Col><TextField label="Initial Step Count" placeholder="0" type="number" value={count} onChange={(e) => setCount(e.target.value)} min={0} />
							</Col>
							<Col>
								<div className="add-user-button-wrapper">
									<Button primary block onClick={handleSubmit}>
										Add User +
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

function StepsIndex({ initialData }) {
	const { data } = useSWR(`${BASE_URL}/steps`, { initialData });

	if (data) {
		// Sort array by total steps.
		var sorted_data = data.data.sort((a, b) => {
			return getTodaysSteps(a.count) < getTodaysSteps(b.count) ? 1 : -1;
		});

		return (
			<>
				<Head>
					<title>User Stats - Pedometer</title>
				</Head>
				<SideNavPanel active_item="steps" />
				<div className="page-content">
					<h1 className="page-title">User Stats</h1>
					<p>
						From this page, the individual user statistics can be browsed. Just click on the user in the list to view the information.
          </p>
					<br />
					<Table bordered>
						<thead>
							<tr>
								<th>#</th>
								<th>User</th>
								<th>
									Step Count
                </th>
							</tr>
						</thead>
						<tbody>
							{
								sorted_data.map((user, index) => (
									<Link key={user.user_name} href={`/steps/${user.user_name}`}>
										<tr className="clickable">
											<td>{index + 1}</td>
											<td>@{user.user_name}</td>
											<td>{getTodaysSteps(user.count)} </td>
										</tr>
									</Link>
								))
							}
							<StepsForm />
						</tbody>
					</Table>
				</div>
			</>
		);
	}

	return <h1>Loading...</h1>;
}

export default StepsIndex;
