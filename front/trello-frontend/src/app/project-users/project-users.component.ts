import { Component, OnInit } from '@angular/core';
import { UserService } from '../service/user.service';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { User } from '../model/user.model';
import { FormsModule } from '@angular/forms';
import { ProjectService } from '../service/project.service';

@Component({
  selector: 'app-project-users',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './project-users.component.html',
  styleUrl: './project-users.component.css'
})
export class ProjectUsersComponent implements OnInit{
  projectId: string | null = null
  members: any[] = [];
  availableUsers: any[] = [];
  selectedUserId: string = '';
  showAddUserForm: boolean = false;

  constructor(private userService: UserService, private route: ActivatedRoute, private projectService: ProjectService) { }
  
  ngOnInit(): void {
    this.projectId = this.route.snapshot.paramMap.get('projectId');
    if (this.projectId) {
      this.loadProjectMembers();
      // this.loadAvailableUsers();
    }
  }

  loadProjectMembers(): void {
    this.userService.getUsersByProjectId(this.projectId!).subscribe({
      next: (data) => {
        this.members = data;
        this.loadAvailableUsers();
        
      },
      error: (error) => {
        console.error('Error fetching project members:', error);
      },
    });
  }

  loadAvailableUsers(): void {
    console.log(this.members)
    this.userService.getUsers().subscribe({
      next: (data) => {
        
        this.availableUsers = data.filter(
          (user: any) =>
            user.role === 'member' && 
            !this.members.some((member) => member.id === user.id) ,
            
        );
      },
      error: (error) => {
        console.error('Error fetching available users:', error);
      },
    });
  }

  addUserToProject(): void {
    if (this.selectedUserId) {
      this.projectService.addUserToProject(this.projectId!, this.selectedUserId).subscribe({
        next: () => {
          console.log(`User ${this.selectedUserId} added to project ${this.projectId}`);
          this.loadProjectMembers(); 
          this.loadAvailableUsers(); 
          this.showAddUserForm = false; 
        },
        error: (error) => {
          console.error('Error adding user to project:', error);
        },
      });
    }
  }

}
