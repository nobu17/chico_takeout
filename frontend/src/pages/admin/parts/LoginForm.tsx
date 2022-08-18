import * as React from "react";
import { Container, Stack, TextField, Button, InputAdornment, IconButton } from "@mui/material";
import { SubmitHandler, useForm } from "react-hook-form";

import { isValidEmail } from "../../../libs/util/Email";
import { RequiredErrorMessage } from "../../../libs/ErrorMessages";
import { Visibility, VisibilityOff } from "@mui/icons-material";

type AdminLoginFormProps = {
  input: LoginInput;
  onSubmit: callbackSubmit;
  onGoogleSubmit?: callbackProviderSubmit;
};
interface callbackSubmit {
  (item: LoginInput): void;
}
interface callbackProviderSubmit {
  (): void;
}

export type LoginInput = {
  email: string;
  password: string;
};

export default function AdminLoginForm(props: AdminLoginFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({ defaultValues: props.input });
  const [showPassword, setShowPassword] = React.useState(false);
  const onSubmit: SubmitHandler<LoginInput> = (data) => {
    props.onSubmit(data);
  };
  function handleClickShowPassword() {
    setShowPassword(!(showPassword))
  }
  const handleEmailValidation = (email: string) => {
    console.log("handleEmailValidation");
    const isValid = isValidEmail(email);
    if (!isValid) {
      return "正しい形式のメールアドレスを入力してください。";
    }
    return;
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="ID"
            fullWidth
            {...register("email", {
              required: { value: true, message: RequiredErrorMessage },
              validate: handleEmailValidation,
            })}
            error={Boolean(errors.email)}
            helperText={errors.email && errors.email.message}
          />
          <TextField
            label="Password"
            fullWidth
            type={showPassword ? 'text' : 'password'}
            autoComplete="current-password"
            {...register("password", {
              required: { value: true, message: RequiredErrorMessage },
            })}
            InputProps={{
              endAdornment:
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={handleClickShowPassword}
                    onMouseDown={(e) => e.preventDefault()}
                    edge="end"
                  >
                    {showPassword ? <VisibilityOff /> : <Visibility />}
                  </IconButton>
                </InputAdornment>
            }}
            error={Boolean(errors.password)}
            helperText={errors.password && errors.password.message}
          />
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            ログイン
          </Button>
          {props.onGoogleSubmit ? (
            <Button
              color="error"
              variant="contained"
              size="large"
              onClick={props.onGoogleSubmit}
            >
              Googleアカウントでログイン
            </Button>
          ) : (
            <></>
          )}
        </Stack>
      </Container>
    </>
  );
}
