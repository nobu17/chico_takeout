import * as React from "react";
import {
  Container,
  Stack,
  TextField,
  Button,
  InputAdornment,
  IconButton,
} from "@mui/material";
import { SubmitHandler, useForm } from "react-hook-form";

import { isValidEmail } from "../../../libs/util/Email";
import { RequiredErrorMessage } from "../../../libs/ErrorMessages";
import { Visibility, VisibilityOff } from "@mui/icons-material";

type SignUpFormProps = {
  input: SignUpInput;
  onSubmit: callbackSubmit;
};
interface callbackSubmit {
  (data: SignUpInput): void;
}

export type SignUpInput = {
  email: string;
  password: string;
};

export default function SignUpForm(props: SignUpFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpInput>({ defaultValues: props.input });
  const [showPassword, setShowPassword] = React.useState(false);
  const onSubmit: SubmitHandler<SignUpInput> = (data) => {
    props.onSubmit(data);
  };
  function handleClickShowPassword() {
    setShowPassword(!showPassword);
  }
  const handleEmailValidation = (email: string) => {
    const isValid = isValidEmail(email);
    if (!isValid) {
      return "正しい形式のメールアドレスを入力してください。";
    }
    return;
  };

  const reg = new RegExp("^(?=.*[A-Z])(?=.*[.?/-])[a-zA-Z0-9.?/-]{8,24}$");
  // const passwordReget = /^(?=.*?[a-z])(?=.*?\d)(?=.*?[!-\/:-@[-`{-~])[!-~]{8,100}$/i;
  const handlePasswordValidation = (password: string) => {
    if (!reg.test(password)) {
      return "英数字の大小文字と記号を含む,8文字以上24文字以下の値を入力してください。";
    }
    return;
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="email"
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
            type={showPassword ? "text" : "password"}
            autoComplete="current-password"
            {...register("password", {
              required: { value: true, message: RequiredErrorMessage },
              validate: handlePasswordValidation,
            })}
            InputProps={{
              endAdornment: (
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
              ),
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
            登録
          </Button>
        </Stack>
      </Container>
    </>
  );
}
