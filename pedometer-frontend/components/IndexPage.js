import Head from "next/head";
import { useRouter } from "next/router";
import React from "react";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import { AiFillGithub } from "react-icons/ai";
import { FaAws, FaReact } from "react-icons/fa";
import { HiArrowRight } from "react-icons/hi";
import { IoHardwareChipOutline } from "react-icons/io5";
import useSWR from "swr";
import { BASE_URL } from "../utils/fetcher";
import SideNavPanel from "./SideNavPanel";



function IndexPage({ initialSteps, initialChallenges }) {
	const router = useRouter();

	const { data: stepsData } = useSWR(`${BASE_URL}/steps`, { initialData: initialSteps });
	const { data: challengesData } = useSWR(`${BASE_URL}/challenges`, { initialData: initialChallenges });

	const defaultChallenge = challengesData.data && challengesData.data.find((val) => val.challenge_name == "default");

	return <>
		<Head>
			<title>Dashboard - Pedometer</title>
		</Head>
		<SideNavPanel active_item="index" />
		<div className="page-content">
			<Container>
				<h1 className="page-title">Dashboard</h1>
				<br></br>
				<br></br>
				<Row>
					<Col>
						<Card>
							<Card.Header>Users</Card.Header>
							<Card.Body className="center-text">
								<Card.Text>
									<div className="quick-stat-number">{stepsData.data ? stepsData.data.length : 0}</div>
                  users registererd
                </Card.Text>
								<br></br>
								<Button onClick={() => router.push("/steps")} variant="outline-info">
									View all users stats {" "}
									<HiArrowRight />
								</Button>
							</Card.Body>
						</Card>
						<br></br>
					</Col>
					<Col>
						<Card>
							<Card.Header>Challenges</Card.Header>
							<Card.Body className="center-text">
								<Card.Text>
									<div className="quick-stat-number">{challengesData.data ? challengesData.data.length : 0}</div>
                  challenges started
                </Card.Text>
								<br></br>
								<Button onClick={() => router.push("/challenges")} variant="outline-info">
									View all challenges  {" "}
									<HiArrowRight />
								</Button>
							</Card.Body>
						</Card>            <br></br>
					</Col>
					<Col>
						<Card>
							<Card.Header>Steps</Card.Header>
							<Card.Body className="center-text">
								<Card.Text>
									<div className="quick-stat-number">{defaultChallenge && defaultChallenge.current}</div>
                  total steps walked
                </Card.Text>
								<br></br>
								<Button onClick={() => router.push("/steps")} variant="outline-info">
									View steps by user  {" "}
									<HiArrowRight />
								</Button>
							</Card.Body>
						</Card>
						<br></br>
					</Col>
				</Row>

				<hr></hr>

				<h3 className="page-title">About this project</h3>
				<p>
					This project was done by a team of EIE students from Imperial College London students as part of their ELEC50009 Information Processing module.
          The members were{" "}
					<a href="https://github.com/jjlehner"><AiFillGithub />Jonah Lehner</a>,{" "}
					<a href="https://github.com/SamtheSaint"><AiFillGithub />Samuel Adekunle</a>,{" "}
					<a href="https://github.com/neeldug"><AiFillGithub />Neel Dugar</a>,{" "}
					<a href="https://github.com/Sam-R-Taylor"><AiFillGithub />Sam Taylor</a>,{" "}
					<a href="https://github.com/max-wickham"><AiFillGithub />Maximus Wickham</a> and {" "}
					<a href="https://github.com/timeo-schmidt"><AiFillGithub />Timeo Schmidt</a>.
          <br></br>
					<br></br>
          To find out more about the project, check out the relevant GitHub repositories here: {" "}
					<a href="https://github.com/SamtheSaint/pedometer-api"><FaAws /> AWS API</a>, {" "}
					<a href="https://github.com/SamtheSaint/pedometer-frontend"><FaReact /> React Frontend</a>, {" "}
					<a href=""><IoHardwareChipOutline /> Hardware Design</a> {" "}.
          <br></br>
					<br></br>
          Thank You!
        </p>
			</Container>

			<br />
		</div>
	</>;
}

export default IndexPage;
