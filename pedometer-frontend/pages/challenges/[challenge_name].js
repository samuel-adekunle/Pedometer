import ChallengesPage from "../../components/ChallengesPage";
import { BASE_URL, fetcher } from "../../utils/fetcher";

function Challenge({ initialData, challengeName }) {
	return <ChallengesPage initialData={initialData} challengeName={challengeName} />;
}

export async function getServerSideProps({ params }) {
	const { challenge_name: challengeName } = params;
	const initialData = await fetcher(`${BASE_URL}/challenges/${challengeName}`);
	return { props: { initialData, challengeName } };
}

export default Challenge;