import { Backdrop, CircularProgress, Typography } from "@mui/material";

type LoadingSpinnerProps = {
  isLoading: boolean;
  message?: string;
};

export default function LoadingSpinner(props: LoadingSpinnerProps) {
  return (
    <>
      <Backdrop open={props.isLoading}>
        <CircularProgress color="warning" size={70} />
        <Typography color="text.primary" position="absolute" sx={{ mt: 13 }}>
          {props.message}
        </Typography>
      </Backdrop>
    </>
  );
}
