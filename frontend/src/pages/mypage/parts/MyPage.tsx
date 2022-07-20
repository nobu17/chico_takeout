import PageTitle from "../../../components/parts/PageTitle";
import PageHeader from "../../../components/parts/PageHeader";

type MyPageProps = {
  title: string;
  links?: PageLink[];
  children?: JSX.Element;
};

type PageLink = {
  title: string;
  url: string;
};

const defaultLinks = [{ title: "マイページ", url: "/my_page" }];

export default function MyPage(props: MyPageProps) {
  const link = props.links ? props.links : defaultLinks;
  return (
    <>
      <PageHeader currentTitle={props.title} links={link} />
      <PageTitle title={props.title} />
      {props.children}
    </>
  );
}
