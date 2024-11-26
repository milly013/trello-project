import { RouterModule, Routes } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { AddTaskComponent } from './add-task/add-task.component';
import { AddUserComponent } from './add-user/add-user.component';
import { RemoveUserComponent } from './remove-user/remove-user.component';
import { RegistrationComponent } from './registration/registration.component';
import { AddMemberToTaskComponent } from './add-member-to-task/add-member-to-task.component';
import { RemoveMemberFromTaskComponent } from './remove-member-from-task/remove-member-from-task.component';
import { VerificationComponent } from './verification/verification.component';
import { HomeComponentComponent } from './home-component/home-component.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { UserListComponent } from './user-list/user-list.component';
import { LoginComponent } from './login/login.component';
import { ProjectDetailsComponent } from './project-details/project-details.component';
import { TaskStatusComponent } from './task-status/task-status.component';

import { ForgotPasswordComponent } from './forgot-password/forgot-password.component';
import { ResetPasswordComponent } from './reset-password/reset-password.component';

import { TaskListComponent } from './task-list/task-list.component';
import { TaskUsersComponent } from './task-users/task-users.component';
import { ChangePasswordComponent } from './change-password/change-password.component';
import { UserProfileComponent } from './user-profile/user-profile.component';


export const appRoutes: Routes = [
  { path: 'project-list/create-project', component: CreateProjectComponent },
  { path: 'add-task/:projectId',component: AddTaskComponent},
  { path:'add-user',component:AddUserComponent},
  { path:'remove-user'  ,component:RemoveUserComponent},
  { path:'home-page/registration',component:RegistrationComponent},
  { path:'verification',component:VerificationComponent},
  { path:'add-member-to-task',component:AddMemberToTaskComponent},
  { path: 'remove-member-from-task', component: RemoveMemberFromTaskComponent },
  { path: 'home-page', component: HomeComponentComponent},
  { path: 'project-list', component: ProjectListComponent},
  { path: 'user-list', component: UserListComponent},
  { path: 'home-page/login', component: LoginComponent},
  { path: 'project-details/:id', component: ProjectDetailsComponent },
  { path: 'task-status',component:TaskStatusComponent},
  { path: 'home-page/login', component: LoginComponent},
  { path: 'project-details/:id', component: ProjectDetailsComponent },

  { path: 'forgot-password', component: ForgotPasswordComponent },
  { path: 'reset-password', component: ResetPasswordComponent },

  { path: 'task-list/:projectId', component: TaskListComponent },
  { path: 'task-users/:taskId', component: TaskUsersComponent },
  { path: 'change-password', component: ChangePasswordComponent },
  {path: 'user-profile',component: UserProfileComponent}



];
 