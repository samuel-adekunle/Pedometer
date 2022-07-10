import StepsPage from "../../components/StepsPage";
import { BASE_URL, fetcher } from "../../utils/fetcher";

function Step({ initialData, userName }) {
	return <StepsPage initialData={initialData} userName={userName} />;
}

export async function getServerSideProps({ params }) {
	const { user_name: userName } = params;
	const initialData = await fetcher(`${BASE_URL}/steps/${userName}`);
	return { props: { initialData, userName } };
}

export default Step;