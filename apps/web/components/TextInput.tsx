import { cn } from "@linkwarden/lib/utils";
import React, { forwardRef } from "react";

export type TextInputProps = React.ComponentPropsWithoutRef<"input">;

const TextInput = forwardRef<HTMLInputElement, TextInputProps>(
  ({ className, type = "text", ...rest }, ref) => {
    return (
      <input
        ref={ref}
        type={type}
        className={cn(
          "w-full rounded-md p-2 border-neutral-content border-solid border outline-none focus:border-primary focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:ring-offset-1 ring-offset-background duration-100",
          className
        )}
        {...rest}
      />
    );
  }
);

TextInput.displayName = "TextInput";

export default TextInput;
