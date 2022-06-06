import PageTitle from "./PageTitle";
import PageHeader from "./PageHeader";

type AdminPageProps = {
  title: string;
  links?: PageLink[];
  children?: JSX.Element;
};

type PageLink = {
  title: string;
  url: string;
};

const defaultLinks = [{ title: "管理ページ", url: "/admin" }];

export default function AdminPage(props: AdminPageProps) {
  const link = props.links ? props.links : defaultLinks;
  return (
    <>
      <PageHeader currentTitle={props.title} links={link} />
      <PageTitle title={props.title} />
      {props.children}
    </>
  );
}
