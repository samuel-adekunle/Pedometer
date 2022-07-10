import { BASE_URL, fetcher } from "../../utils/fetcher";
import StepsIndex from "../../components/StepsIndex";

function Steps({ initialData }) {
	return <StepsIndex initialData={initialData} />;
}

export async function getServerSideProps() {
	const initialData = await fetcher(`${BASE_URL}/steps`);
	return { props: { initialData } };
}

export default Steps;