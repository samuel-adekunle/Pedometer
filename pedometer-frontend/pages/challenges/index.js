import ChallengesIndex from "../../components/ChallengesIndex";
import { BASE_URL, fetcher } from "../../utils/fetcher";

function Challenges({ initialData }) {
	return <ChallengesIndex initialData={initialData} />;
}

export async function getServerSideProps() {
	const initialData = await fetcher(`${BASE_URL}/challenges`);
	return { props: { initialData } };
}

export default Challenges;