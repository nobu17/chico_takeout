import Typography from "@mui/material/Typography";

type PageTitleProps = {
  title: string;
};

export default function PageTitle(props: PageTitleProps) {
  return (
    <>
      <Typography
        component="h2"
        variant="h3"
        align="center"
        color="text.primary"
        gutterBottom
      >
        {props.title}
      </Typography>
    </>
  );
}
