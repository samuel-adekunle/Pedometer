import IndexPage from "../components/IndexPage";
import { BASE_URL, fetcher } from "../utils/fetcher";

function Index({ initialSteps, initialChallenges }) {
	return <IndexPage initialSteps={initialSteps} initialChallenges={initialChallenges} />;
}

export async function getServerSideProps() {
	const initialSteps = await fetcher(`${BASE_URL}/steps`);
	const initialChallenges = await fetcher(`${BASE_URL}/challenges`);
	return { props: { initialSteps, initialChallenges } };
}

export default Index;
