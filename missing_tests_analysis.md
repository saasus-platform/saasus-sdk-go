# 不足しているテストメソッド一覧

## インタフェースメソッド全169個と既存テスト84個の比較

### 不足しているテストメソッド（85個）

#### 1. WithBodyWithResponse系メソッド（42個）
1. `UpdateAuthInfoWithBodyWithResponse`
2. `LinkAwsMarketplaceWithBodyWithResponse`
3. `SignUpWithAwsMarketplaceWithBodyWithResponse`
4. `ConfirmSignUpWithAwsMarketplaceWithBodyWithResponse`
5. `UpdateBasicInfoWithBodyWithResponse`
6. `CreateAuthCredentialsWithBodyWithResponse`
7. `UpdateCustomizePageSettingsWithBodyWithResponse`
8. `UpdateCustomizePagesWithBodyWithResponse`
9. `CreateEnvWithBodyWithResponse`
10. `UpdateEnvWithBodyWithResponse`
11. `ConfirmExternalUserLinkWithBodyWithResponse`
12. `RequestExternalUserLinkWithBodyWithResponse`
13. `UpdateIdentityProviderWithBodyWithResponse`
14. `ValidateInvitationWithBodyWithResponse`
15. `UpdateNotificationMessagesWithBodyWithResponse`
16. `CreateRoleWithBodyWithResponse`
17. `CreateSaasUserAttributeWithBodyWithResponse`
18. `UpdateSignInSettingsWithBodyWithResponse`
19. `SignUpWithBodyWithResponse`
20. `ResendSignUpConfirmationEmailWithBodyWithResponse`
21. `UpdateSingleTenantSettingsWithBodyWithResponse`
22. `CreateTenantAttributeWithBodyWithResponse`
23. `CreateTenantWithBodyWithResponse`
24. `UpdateTenantWithBodyWithResponse`
25. `UpdateTenantBillingInfoWithBodyWithResponse`
26. `UpdateTenantIdentityProviderWithBodyWithResponse`
27. `CreateTenantInvitationWithBodyWithResponse`
28. `UpdateTenantPlanWithBodyWithResponse`
29. `CreateTenantUserWithBodyWithResponse`
30. `UpdateTenantUserWithBodyWithResponse`
31. `CreateTenantUserRolesWithBodyWithResponse`
32. `CreateUserAttributeWithBodyWithResponse`
33. `CreateSaasUserWithBodyWithResponse`
34. `UpdateSaasUserAttributesWithBodyWithResponse`
35. `UpdateSaasUserEmailWithBodyWithResponse`
36. `ConfirmEmailUpdateWithBodyWithResponse`
37. `RequestEmailUpdateWithBodyWithResponse`
38. `UpdateUserMfaPreferenceWithBodyWithResponse`
39. `UpdateSoftwareTokenWithBodyWithResponse`
40. `CreateSecretCodeWithBodyWithResponse`
41. `UpdateSaasUserPasswordWithBodyWithResponse`

#### 2. 特殊なメソッド（2個）
42. `ReturnInternalServerErrorWithResponse`
43. `ResendSignUpConfirmationEmailWithResponse`

### 既存のテストメソッド（84個）
1. `TestE2E_GetAuthInfoWithResponse`
2. `TestE2E_UpdateAuthInfoWithResponse`
3. `TestE2E_GetBasicInfoWithResponse`
4. `TestE2E_UpdateBasicInfoWithResponse`
5. `TestE2E_GetSaasUsersWithResponse`
6. `TestE2E_CreateSaasUserWithResponse`
7. `TestE2E_GetSaasUserWithResponse`
8. `TestE2E_DeleteSaasUserWithResponse`
9. `TestE2E_GetTenantsWithResponse`
10. `TestE2E_CreateTenantWithResponse`
11. `TestE2E_GetRolesWithResponse`
12. `TestE2E_CreateRoleWithResponse`
13. `TestE2E_GetEnvsWithResponse`
14. `TestE2E_GetUserAttributesWithResponse`
15. `TestE2E_GetTenantAttributesWithResponse`
16. `TestE2E_GetSignInSettingsWithResponse`
17. `TestE2E_GetIdentityProvidersWithResponse`
18. `TestE2E_GetCustomizePagesWithResponse`
19. `TestE2E_GetCustomizePageSettingsWithResponse`
20. `TestE2E_FindNotificationMessagesWithResponse`
21. `TestE2E_GetSingleTenantSettingsWithResponse`
22. `TestE2E_GetCloudFormationLaunchStackLinkForSingleTenantWithResponse`
23. `TestE2E_ResetPlanWithResponse`
24. `TestE2E_CreateTenantAndPricingWithResponse`
25. `TestE2E_DeleteStripeTenantAndPricingWithResponse`
26. `TestE2E_SignUpWithResponse`
27. `TestE2E_GetTenantWithResponse`
28. `TestE2E_DeleteTenantWithResponse`
29. `TestE2E_GetTenantUsersWithResponse`
30. `TestE2E_CreateTenantUserWithResponse`
31. `TestE2E_GetRoleWithResponse`
32. `TestE2E_DeleteRoleWithResponse`
33. `TestE2E_CreateUserAttributeWithResponse`
34. `TestE2E_DeleteUserAttributeWithResponse`
35. `TestE2E_CreateTenantAttributeWithResponse`
36. `TestE2E_DeleteTenantAttributeWithResponse`
37. `TestE2E_CreateEnvWithResponse`
38. `TestE2E_GetEnvWithResponse`
39. `TestE2E_DeleteEnvWithResponse`
40. `TestE2E_UpdateSignInSettingsWithResponse`
41. `TestE2E_UpdateCustomizePagesWithResponse`
42. `TestE2E_UpdateCustomizePageSettingsWithResponse`
43. `TestE2E_UpdateNotificationMessagesWithResponse`
44. `TestE2E_UpdateSingleTenantSettingsWithResponse`
45. `TestE2E_GetAuthCredentialsWithResponse`
46. `TestE2E_CreateAuthCredentialsWithResponse`
47. `TestE2E_UpdateEnvWithResponse`
48. `TestE2E_LinkAwsMarketplaceWithResponse`
49. `TestE2E_SignUpWithAwsMarketplaceWithResponse`
50. `TestE2E_ConfirmSignUpWithAwsMarketplaceWithResponse`
51. `TestE2E_RequestExternalUserLinkWithResponse`
52. `TestE2E_ConfirmExternalUserLinkWithResponse`
53. `TestE2E_UpdateIdentityProviderWithResponse`
54. `TestE2E_GetAllTenantUsersWithResponse`
55. `TestE2E_GetAllTenantUserWithResponse`
56. `TestE2E_UpdateTenantWithResponse`
57. `TestE2E_GetTenantUserWithResponse`
58. `TestE2E_UpdateTenantUserWithResponse`
59. `TestE2E_DeleteTenantUserWithResponse`
60. `TestE2E_ValidateInvitationWithResponse`
61. `TestE2E_GetInvitationValidityWithResponse`
62. `TestE2E_UpdateTenantBillingInfoWithResponse`
63. `TestE2E_GetTenantIdentityProvidersWithResponse`
64. `TestE2E_UpdateTenantIdentityProviderWithResponse`
65. `TestE2E_GetTenantInvitationsWithResponse`
66. `TestE2E_CreateTenantInvitationWithResponse`
67. `TestE2E_GetTenantInvitationWithResponse`
68. `TestE2E_DeleteTenantInvitationWithResponse`
69. `TestE2E_UpdateTenantPlanWithResponse`
70. `TestE2E_GetStripeCustomerWithResponse`
71. `TestE2E_CreateTenantUserRolesWithResponse`
72. `TestE2E_DeleteTenantUserRoleWithResponse`
73. `TestE2E_GetUserInfoWithResponse`
74. `TestE2E_UpdateSaasUserAttributesWithResponse`
75. `TestE2E_UpdateSaasUserEmailWithResponse`
76. `TestE2E_ConfirmEmailUpdateWithResponse`
77. `TestE2E_RequestEmailUpdateWithResponse`
78. `TestE2E_GetUserMfaPreferenceWithResponse`
79. `TestE2E_UpdateUserMfaPreferenceWithResponse`
80. `TestE2E_UpdateSoftwareTokenWithResponse`
81. `TestE2E_CreateSecretCodeWithResponse`
82. `TestE2E_UpdateSaasUserPasswordWithResponse`
83. `TestE2E_UnlinkProviderWithResponse`
84. `TestE2E_CreateSaasUserAttributeWithResponse`

## 分析結果

### 主な不足パターン
1. **WithBodyWithResponse系**: 多くのメソッドで`WithBodyWithResponse`バリアントのテストが不足
2. **エラーテスト**: `ReturnInternalServerErrorWithResponse`のテストが不足
3. **メール関連**: `ResendSignUpConfirmationEmailWithResponse`のテストが不足

### 優先度
1. **高**: `ReturnInternalServerErrorWithResponse` - エラーハンドリングテスト
2. **高**: `ResendSignUpConfirmationEmailWithResponse` - 重要な機能
3. **中**: WithBodyWithResponse系 - 代替実装のテスト