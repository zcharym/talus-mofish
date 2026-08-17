export async function sendMagicLinkEmail(
  apiKey: string,
  fromEmail: string,
  toEmail: string,
  appName: string,
  verifyURL: string,
): Promise<void> {
  const response = await fetch('https://api.resend.com/emails', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${apiKey}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      from: fromEmail,
      to: [toEmail],
      subject: `Sign in to ${appName}`,
      html: `<!DOCTYPE html>
<html lang="en">
<body style="font-family: system-ui, sans-serif; padding: 24px;">
  <h1>Sign in to ${appName}</h1>
  <p>Click the link below to finish signing in. This link expires in 10 minutes.</p>
  <p><a href="${verifyURL}">Sign in</a></p>
  <p style="color: #666; font-size: 14px;">If you did not request this email, you can ignore it.</p>
</body>
</html>`,
    }),
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`resend API error (${response.status}): ${body}`);
  }
}
