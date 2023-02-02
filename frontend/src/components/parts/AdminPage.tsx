import PageTitle from "./PageTitle";
import PageHeader from "./PageHeader";
import Container from "@mui/material/Container";

import { useReloadTimer } from "../../hooks/UseTimer";

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
  useReloadTimer(30);
  const link = props.links ? props.links : defaultLinks;
  return (
    <Container>
      <PageHeader currentTitle={props.title} links={link} />
      <PageTitle title={props.title} />
      {props.children}
    </Container>
  );
}
