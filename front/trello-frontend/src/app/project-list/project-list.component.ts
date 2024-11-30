import { Component, OnInit } from '@angular/core';
import { Project, ProjectService } from '../service/project.service';
import { nextTick } from 'process';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterOutlet } from '@angular/router';
import { AuthService } from '../service/auth.service';


@Component({
  selector: 'app-project-list',
  standalone: true,
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.css'],
  imports: [CommonModule,RouterLink,RouterOutlet,]
})
export class ProjectListComponent implements OnInit {
  projects: Project[] = [];

  constructor(private projectService: ProjectService, private authService: AuthService) {}

  ngOnInit(): void {
    const userRole = this.authService.getUserRole();
    const userId = this.authService.getUserId();

    if (userRole && userId) {
      if (userRole === 'manager') {
        this.getProjectsByManager(userId);
      } else if (userRole === 'member') {
        this.getProjectsByMember(userId);
      }
    } else {
      console.error('User information is missing or token is invalid');
    }
  }

  getProjectsByManager(managerId: string): void {
    this.projectService.getProjectsByManager(managerId).subscribe({
      next: (data: Project[]) => {
        this.projects = data;
      },
      error: (error: any) => {
        console.error('Error fetching projects for manager', error);
      },
      complete: () => {
        console.log('Fetching projects for manager complete');
      }
    });
  }

  getProjectsByMember(memberId: string): void {
    this.projectService.getProjectsByMember(memberId).subscribe({
      next: (data: Project[]) => {
        this.projects = data;
      },
      error: (error: any) => {
        console.error('Error fetching projects for member', error);
      },
      complete: () => {
        console.log('Fetching projects for member complete');
      }
    });
  }
  isUserManager():boolean{
    return this.authService.getUserRole() === 'manager';
  }
  
}
